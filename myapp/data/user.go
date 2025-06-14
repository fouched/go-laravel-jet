package data

import (
	"errors"
	"github.com/fouched/celeritas"
	up "github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type User struct {
	ID        int       `db:"id,omitempty"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Active    int       `db:"user_active"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Token     Token     `db:"-"`
}

// Table allows specifying a custom table name for the model
func (u *User) Table() string {
	return "users"
}

func (u *User) Validate(validator *celeritas.Validation) {
	validator.Check(strings.TrimSpace(u.LastName) != "", "last_name", "Last name is required")
	validator.Check(strings.TrimSpace(u.FirstName) != "", "first_name", "First name is required")
	validator.IsEmail(celeritas.Field{
		Name:  "email",
		Label: "Email",
		Value: strings.TrimSpace(u.Email),
	})
}

func (u *User) GetAll() ([]*User, error) {
	var users []*User

	collection := upper.Collection(u.Table())
	rs := collection.Find().OrderBy("last_name")
	err := rs.All(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *User) GetByEmail(email string) (*User, error) {
	var user User

	collection := upper.Collection(u.Table())
	rs := collection.Find(up.Cond{"email =": email}) // = is assumed, use it for clarity
	err := rs.One(&user)
	if err != nil {
		return nil, err
	}

	err = setUserToken(&user)
	if err != nil {
		if !errors.Is(err, up.ErrNilRecord) && !errors.Is(err, up.ErrNoMoreRows) {
			return nil, err
		}
	}

	return &user, nil
}

func (u *User) Get(id int) (*User, error) {
	var user User

	collection := upper.Collection(u.Table())
	rs := collection.Find(up.Cond{"id =": id}) // = is assumed, use it for clarity
	err := rs.One(&user)
	if err != nil {
		return nil, err
	}

	err = setUserToken(&user)
	if err != nil {
		if !errors.Is(err, up.ErrNilRecord) && !errors.Is(err, up.ErrNoMoreRows) {
			return nil, err
		}
	}

	return &user, nil
}

func setUserToken(user *User) error {
	var token Token
	collection := upper.Collection(token.Table())
	rs := collection.Find(
		up.Cond{
			"user_id =": user.ID,
			"expiry >":  time.Now(),
		}).OrderBy("created_at desc")

	err := rs.One(&token)
	if err != nil {
		return err
	}

	user.Token = token

	return nil
}

func (u *User) Update(user User) error {
	user.UpdatedAt = time.Now()
	collection := upper.Collection(u.Table())
	rs := collection.Find(user.ID)

	err := rs.Update(&user)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Delete(id int) error {
	collection := upper.Collection(u.Table())
	rs := collection.Find(id)
	err := rs.Delete()
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Insert(user User) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Password = string(hash)

	collection := upper.Collection(u.Table())
	rs, err := collection.Insert(user)
	if err != nil {
		return 0, err
	}

	id := getInsertID(rs.ID())
	return id, nil
}

func (u *User) ResetPassword(id int, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	user, err := u.Get(id)
	if err != nil {
		return err
	}

	u.Password = string(hash)

	err = user.Update(*u)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (u *User) CheckForRememberToken(id int, token string) bool {
	var rememberToken RememberToken
	rt := RememberToken{}
	collection := upper.Collection(rt.Table())
	rs := collection.Find(up.Cond{"user_id": id, "remember_token": token})
	err := rs.One(&rememberToken)
	return err == nil
}
