package plain

import "github.com/rugwirobaker/paypack-backend/pkg/errors"

// Compare plain passwords
func Compare(a, b string) error {
	const op errors.Op = "pkg/passwords/plain/Compare"

	if a != b {
		return errors.E(op, "invalid login data: wrong password", errors.KindBadRequest)
	}
	return nil
}
