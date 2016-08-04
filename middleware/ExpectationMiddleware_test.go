package middleware

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("ExpectationMiddleware", func() {
    var em *ExpectationMiddleware

    BeforeEach(func() {
        em = &ExpectationMiddleware{}
    })

    Describe("HandleBuffer", func() {

    })

    Describe("AddExpectation", func() {
        It("Now contains the expectation", func() {
            eh := func([]byte) ([]byte, error) {
                return nil, nil
            }

            em.AddExpectation(eh)
            Expect(em.expectations).To(ConsistOf(eh))
        }
    })

    Describe("HasFailedExpectations", func() {
        Context("Has failed expectations", func() {
            BeforeEach(func() {
                em.errs = []error{fmt.Errorf("A thing happened!")}
            })

            It("Reports an error", func() {
                Expect(em.HasFailedExpectations()).To(BeTrue())
            })
        })

        Context("Has no failed expectations", func() {
            It("Does not report an error", func() {
                Expect(em.HasFailedExpectations()).To(BeFalse())
            })
        })
    })

    Describe("HasRemainingExpectations", func() {
        Context("Has remaining expectations", func() {
            BeforeEach(func() {
                em.expectations = []ExpectationHandler{func([]byte) ([]byte, error) {
                    return nil, nil
                }}
            })

            It("Reports true", func() {
                Expect(em.HasRemainingExpectations()).To(BeTrue())
            })
        })

        Context("Has no remaining expectations", func() {
            It("Does not report true", func() {
                Expect(em.HasRemainingExpectations()).To(BeFalse())
            })
        })
    })
})