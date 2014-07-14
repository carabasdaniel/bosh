package disk_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh/platform/disk"
	fakesys "bosh/system/fakes"
	"fmt"
)

func init() {
	Describe("Testing with Ginkgo", func() {
		It("windows mount searcher test", func() {
			fakeFs := fakesys.NewFakeFileSystem()

			searcher := NewWindowsMountsSearcher(fakeFs)
			mounts, err := searcher.SearchMounts()
			fmt.Println("-----------------------------------------------")
			for _, k := range mounts {
				fmt.Println(k.MountPoint + " -> " + k.PartitionPath)
			}
			fmt.Println("-----------------------------------------------")
			Expect(err).ToNot(HaveOccurred())
		})
	})
}
