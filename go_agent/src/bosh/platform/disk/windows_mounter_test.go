package disk_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh/platform/disk"
	"time"

	fakedisk "bosh/platform/disk/fakes"
	fakesys "bosh/system/fakes"
)

func init() {
	Describe("mounterTest", func() {
		It("windows mount test", func() {
			fakeRunner := fakesys.NewFakeCmdRunner()
			fakeFs := fakesys.NewFakeFileSystem()

			searcher := NewWindowsMountsSearcher(fakeFs)
			mounter := NewFakeWindowsMounter(fakeRunner, searcher, 1*time.Second, fakedisk.NewFakeDiskPart())
			err := mounter.Mount("2", "C:\\mountP")
			Expect(err).ToNot(HaveOccurred())
		})
	})
	Describe("isMountedTest", func() {
		It("windows is mount test", func() {
			fakeRunner := fakesys.NewFakeCmdRunner()
			fakeFs := fakesys.NewFakeFileSystem()

			searcher := NewWindowsMountsSearcher(fakeFs)

			mounter := NewFakeWindowsMounter(fakeRunner, searcher, 1*time.Second, fakedisk.NewFakeDiskPart())
			ok, err := mounter.IsMounted("C:\\mountP\\")
			Expect(ok).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())
		})
	})
}
