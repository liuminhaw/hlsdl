/*
Copyright © 2023 Min-Haw, Liu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"

	"github.com/liuminhaw/hlsdl"
	"github.com/spf13/cobra"
)

// segmentsCmd represents the segments command
var segmentsCmd = &cobra.Command{
	Use:   "segments <m3u8-url>",
	Short: "Parse and show segments information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		m3u8Url := args[0]

		segments, err := hlsdl.ParseSegments(m3u8Url, map[string]string{})
		if err != nil {
			return err
		}

		for _, segment := range segments {
			fmt.Printf("Sequence id: %d\n", segment.SeqId)
			fmt.Printf("Title: %s\n", segment.Title)
			fmt.Printf("URI: %s\n", segment.URI)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(segmentsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// segmentsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// segmentsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
