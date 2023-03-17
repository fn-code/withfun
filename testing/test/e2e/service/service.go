package service

import (
	"github.com/fn-code/withfun/testing/internal/service"
	"github.com/fn-code/withfun/testing/internal/storage"
	strgmock "github.com/fn-code/withfun/testing/internal/storage/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service Sum", func() {

	var (
		mockCtrl    *gomock.Controller //gomock struct
		storageMock *strgmock.MockStorage
		svcImpl     *service.Impl
	)

	//initialization
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		storageMock = strgmock.NewMockStorage(mockCtrl)
		svcImpl = &service.Impl{
			Storage: storageMock,
		}
	})

	//tear down
	AfterEach(func() {
		// We add this in order to check that all the registered mocks ware really called
		mockCtrl.Finish()
	})

	Context("Custom sum service", func() {

		When("service passed", func() {
			BeforeEach(func() {
				storageMock.EXPECT().AddCustom(&storage.Custom{
					A:         1,
					B:         1,
					OnProcess: nil,
				}).Return(2, nil)
			})
			It("should be success", func() {
				res, err := svcImpl.SumCustom(1, 1, nil)
				Expect(err).To(BeNil())
				Expect(res).To(Equal(2))
			})
		})

		When("service return error", func() {
			BeforeEach(func() {
				storageMock.EXPECT().AddCustom(&storage.Custom{
					A:         1,
					B:         1,
					OnProcess: nil,
				}).Return(2, storage.ErrZoro)
			})
			It("should be fail", func() {

				res, err := svcImpl.SumCustom(1, 1, nil)
				Expect(err).To(Equal(storage.ErrZoro))
				Expect(res).To(Equal(2))
			})
		})

	})
})
