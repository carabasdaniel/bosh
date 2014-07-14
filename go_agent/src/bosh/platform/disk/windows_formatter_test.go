package disk_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh/platform/disk"
	fakedisk "bosh/platform/disk/fakes"
	fakesys "bosh/system/fakes"
)

func init() {
	Describe("Testing with Ginkgo", func() {
		It("windows format with path", func() {

			fakeRunner := fakesys.NewFakeCmdRunner()
			fakeFs := fakesys.NewFakeFileSystem()

			formatter := NewWindowsFormatter(fakeRunner, fakeFs)
			formatter.Format("D:\\", FileSystemNtfs)

			Expect(1).To(Equal(len(fakeRunner.RunCommands)))
		})
		It("windows format with volume id", func() {
			fakeRunner := fakesys.NewFakeCmdRunner()
			fakeFs := fakesys.NewFakeFileSystem()

			formatter := NewWindowsFormatter(fakeRunner, fakeFs)
			formatter.ChangeDiskPart(fakedisk.NewFakeDiskPart())

			err := formatter.Format("1", FileSystemNtfs)
			Expect(err).ToNot(HaveOccurred())
		})
	})
}
