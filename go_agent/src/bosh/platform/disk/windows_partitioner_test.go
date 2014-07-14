package disk_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	boshlog "bosh/logger"
	. "bosh/platform/disk"

	fakedisk "bosh/platform/disk/fakes"
	fakesys "bosh/system/fakes"
)

func init() {
	Describe("Testing with Ginkgo", func() {
		It("windows partitioner", func() {

			fakeRunner := fakesys.NewFakeCmdRunner()

			partitions := []Partition{
				{SizeInMb: 100, Type: PartitionTypeWindows},
				{SizeInMb: 1000, Type: PartitionTypeWindows},
			}

			logger := boshlog.NewLogger(boshlog.LevelNone)
			partitioner := NewFakeWindowsPartitioner(logger, fakeRunner, fakedisk.NewFakeDiskPart())
			err := partitioner.Partition("1", partitions)
			Expect(err).NotTo(HaveOccurred())
		})
		It("Get device remaining size", func() {

			fakeRunner := fakesys.NewFakeCmdRunner()
			logger := boshlog.NewLogger(boshlog.LevelNone)
			partitioner1 := NewFakeWindowsPartitioner(logger, fakeRunner, fakedisk.NewFakeDiskPart())
			result, _ := partitioner1.GetDeviceSizeInMb("1")
			Expect(result).To(Equal(uint64(121)))
		})
	})
}
