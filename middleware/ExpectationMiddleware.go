package middleware

import "fmt"

type ExpectationHandler func([]byte) ([]byte, error)

type ExpectationMiddleware struct {
    expectations []ExpectationHandler
    errs []error
}

func (em *ExpectationMiddleware) HandleBuffer(buff []byte) ([]byte, error) {
    if em.HasRemainingExpectations() {
        next := em.expectations[0]
        em.expectations = em.expectations[1:]
        resp, err := next(buff)
        if err != nil {
            em.errs = append(em.errs, err)
        }

        return resp, nil
    }

    err := fmt.Errorf("Unexpected data on connection!")
    em.errs = append(em.errs, err)
    return nil, err
}

func (em *ExpectationMiddleware) AddExpectation(expectation ExpectationHandler) {
    em.expectations = append(em.expectations, expectation)
}

func (em *ExpectationMiddleware) HasFailedExpectations() bool {
    return len(em.errs) > 0
}

func (em *ExpectationMiddleware) HasRemainingExpectations() bool {
    return len(em.expectations) > 0
}