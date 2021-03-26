// Code generated by vfsgen; DO NOT EDIT.

package rbac

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	pathpkg "path"
	"time"
)

// Root statically implements the virtual filesystem provided to vfsgen.
var Root = func() http.FileSystem {
	fs := vfsgen۰FS{
		"/": &vfsgen۰DirInfo{
			name:    "/",
			modTime: time.Time{},
		},
		"/auth_proxy_role.yaml": &vfsgen۰CompressedFileInfo{
			name:             "auth_proxy_role.yaml",
			modTime:          time.Time{},
			uncompressedSize: 280,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\xce\xbd\x4a\xc5\x40\x10\xc5\xf1\x7e\x9f\x62\x48\x9f\x88\x9d\x6c\x6b\x61\x6f\x61\x23\xb7\x98\xec\x3d\xe0\x98\x64\x67\x99\x99\x8d\x1f\x4f\x2f\xb1\x12\x09\xb7\x3e\x9c\x1f\x7f\x6e\xf2\x02\x73\xd1\x9a\xc9\x66\x2e\x13\xf7\x78\x53\x93\x6f\x0e\xd1\x3a\x2d\x0f\x3e\x89\xde\xed\xf7\x69\x91\x7a\xcd\xf4\xb8\x76\x0f\xd8\xb3\xae\x48\x1b\x82\xaf\x1c\x9c\x13\x51\xe5\x0d\x99\x9a\xe9\xe7\xd7\x68\xc7\x68\x7d\x85\xe7\x34\x12\x37\x79\x32\xed\xcd\x33\xbd\x0e\x07\x8e\x1a\x52\xfe\xea\xc3\x25\x11\x19\x5c\xbb\x95\xe3\x43\x34\x52\xe8\x82\x6a\xd8\x05\x1f\x9e\x88\x76\xd8\xfc\x0b\x14\x03\x07\x86\xcb\x19\xfc\xbf\xfa\xcc\xf5\x3e\xbf\xa3\x04\x97\x02\xf7\x5b\xfe\x4f\x00\x00\x00\xff\xff\x4b\x41\x4e\x63\x18\x01\x00\x00"),
		},
		"/auth_proxy_role_binding.yaml": &vfsgen۰CompressedFileInfo{
			name:             "auth_proxy_role_binding.yaml",
			modTime:          time.Time{},
			uncompressedSize: 257,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x8d\x31\x4e\xc0\x30\x0c\x45\xf7\x9c\xc2\x17\x68\x11\x1b\xca\x06\x0c\xec\x45\x62\x77\x13\x17\x4c\xd3\x38\xb2\x9d\x8a\x72\x7a\x54\xa9\xb0\x80\xd8\x6c\xe9\xbd\xff\xb0\xf1\x0b\xa9\xb1\xd4\x08\x3a\x63\x1a\xb1\xfb\x9b\x28\x7f\xa2\xb3\xd4\x71\xbd\xb3\x91\xe5\x66\xbf\x0d\x2b\xd7\x1c\xe1\xb1\x74\x73\xd2\x49\x0a\x3d\x70\xcd\x5c\x5f\xc3\x46\x8e\x19\x1d\x63\x00\xa8\xb8\x51\x84\xa6\xf2\x71\x0c\x2a\x85\xe6\x8b\x39\xef\x89\x96\x13\xc1\xc6\x4f\x2a\xbd\xfd\x93\x0b\x00\xbf\x6a\x7f\x8c\x07\xeb\xf3\x3b\x25\xb7\x18\x86\x4b\x78\x26\xdd\x39\xd1\x7d\x4a\xd2\xab\xff\x38\x99\x16\xec\xe5\xfb\xb7\x86\x89\x22\xd8\x61\x4e\x5b\xf8\x0a\x00\x00\xff\xff\x38\xc6\x1a\x54\x01\x01\x00\x00"),
		},
		"/auth_proxy_service.yaml": &vfsgen۰CompressedFileInfo{
			name:             "auth_proxy_service.yaml",
			modTime:          time.Time{},
			uncompressedSize: 268,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x8e\xbd\x0a\xc2\x40\x10\x84\xfb\x7b\x8a\x7d\x81\x2b\xc4\x14\x72\x4f\x21\x08\xf6\xeb\x65\x88\x87\xf7\xc7\xee\x12\xf0\xed\x25\x67\xd2\x59\xd8\xed\xec\x0c\xf3\x0d\xf7\x74\x87\x68\x6a\x35\xd0\x7a\x72\xaf\x54\xe7\x40\x37\xc8\x9a\x22\x5c\x81\xf1\xcc\xc6\xc1\x11\x65\x7e\x20\xeb\x76\x11\xc5\x56\x4d\x5a\xf6\x3d\x73\x45\x38\x64\x86\xf8\xc2\x95\x17\x88\x23\xaa\x5c\x7e\x5a\xbe\xc0\x24\x45\xf5\xba\x43\xbe\x51\xed\x1c\x11\x48\xdf\x6a\x28\x4e\x3b\xe2\x86\xea\x4d\x6c\x30\xfd\xde\xf7\x34\xeb\x3a\x36\x6c\x56\xa0\xcb\x34\x9d\x87\x34\x96\x05\x76\x1d\xcf\x23\xa4\xc8\x88\xd6\xe4\xdf\xd1\x9f\x00\x00\x00\xff\xff\x48\xc2\x77\x74\x0c\x01\x00\x00"),
		},
		"/kustomization.yaml": &vfsgen۰CompressedFileInfo{
			name:             "kustomization.yaml",
			modTime:          time.Time{},
			uncompressedSize: 344,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x8f\xb1\x4e\xc5\x30\x0c\x45\xf7\x7e\x85\xa5\x2e\x30\xb4\x19\xd8\xde\xca\x87\x44\x49\xea\xd7\x58\x24\x71\xe5\x38\x94\xf2\xf5\xe8\x45\x80\x5a\xc4\x66\xf9\x9e\xab\xab\x23\x58\xb9\x49\xc0\x7a\x1b\x26\x10\x4e\x38\x1f\x2e\xa7\xef\xdb\x7a\x2a\x0b\x95\xf5\xe7\x97\xd0\x2d\x28\x16\x13\x06\x25\x2e\xf6\xcc\xff\x97\x5d\xfb\x23\xbc\x72\xce\x58\x14\x34\x22\xdc\x39\x25\xde\xa9\xac\xf0\x02\x89\x0a\x56\xa0\x3b\x1c\xdc\x60\x77\x0f\x82\x61\xa1\xea\x7c\xc2\x61\xec\xb8\x6b\x1a\x61\x13\xfe\x38\xe0\x29\xaa\x6e\xf5\x66\xcc\x4a\x1a\x9b\x9f\x03\x67\xe3\xc5\x95\xf0\x69\xde\x9a\xc7\x49\xbc\x0b\x53\x47\x9f\x87\x11\xf6\x48\xa1\x37\x15\x83\xd6\xc7\x82\x80\xc9\xa8\x42\xa1\x02\x96\x65\x63\x2a\x3a\x0f\x53\x5f\xb0\xbd\x66\x2b\xca\x3b\x85\x5f\xb5\x53\x72\x36\xfe\xf3\xbe\xca\x7e\x05\x00\x00\xff\xff\xac\x90\x8c\xf8\x58\x01\x00\x00"),
		},
		"/leader_election_role.yaml": &vfsgen۰CompressedFileInfo{
			name:             "leader_election_role.yaml",
			modTime:          time.Time{},
			uncompressedSize: 419,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x8e\xb1\x4e\x03\x41\x0c\x44\xfb\xfd\x0a\x2b\xd4\x77\x11\x1d\xba\x1f\xa0\xa7\xa0\x77\x76\x87\x64\x95\xbd\xf5\xca\xf6\x1e\x12\x5f\x8f\xee\x72\xa1\x00\x21\x41\xe5\x91\xed\x79\x33\x0f\xd4\xa0\x73\x36\xcb\x52\x8d\x5c\x28\x09\x15\x70\x82\x12\x0a\xa2\x67\xa9\x63\xe0\x96\x5f\xa1\xeb\xcb\x44\x7a\xe2\x38\x72\xf7\x8b\x68\xfe\xe0\xed\x7e\x7d\xb2\x31\xcb\x71\x79\x0c\xd7\x5c\xd3\x44\x2f\x52\x10\x66\x38\x27\x76\x9e\x02\x51\xe5\x19\xd3\x8e\x1d\xee\xd8\x41\xd7\x37\xed\x05\x36\x85\x81\xb8\xe5\x67\x95\xde\x6c\x35\x0c\x74\x38\x04\x22\x85\x49\xd7\x88\x7d\x17\xa5\xbe\xe5\xf3\xcc\xcd\x02\xd1\x02\x3d\xed\xfb\x33\x7c\x9b\x25\xdb\x4d\xbc\xb3\xc7\xcb\xcd\xa2\x60\xc7\x26\x7b\x4b\x77\xd9\xbe\xee\x09\x05\x8e\xff\xc6\x1f\xcd\xd9\xfb\x2f\x2d\x7e\xe4\xfc\x09\x8e\x05\xd5\xbf\x11\xf7\xf2\x9f\x01\x00\x00\xff\xff\x0e\xcb\x95\x88\xa3\x01\x00\x00"),
		},
		"/leader_election_role_binding.yaml": &vfsgen۰CompressedFileInfo{
			name:             "leader_election_role_binding.yaml",
			modTime:          time.Time{},
			uncompressedSize: 263,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x8d\x3d\x8a\xc3\x30\x10\x46\x7b\x9d\x62\x2e\x60\x2f\xdb\x2d\xea\x76\x9b\xed\x1d\x48\x3f\x96\x3e\x27\x13\xcb\x1a\xa1\x1f\x43\x72\xfa\x60\x70\x92\xce\xdd\x0c\xbc\xf7\x3d\x4e\x72\x46\x2e\xa2\xd1\x52\x1e\xd9\xf5\xdc\xea\x55\xb3\x3c\xb8\x8a\xc6\x7e\xfe\x29\xbd\xe8\xd7\xfa\x6d\x66\x89\xde\xd2\xa0\x01\x7f\x12\xbd\xc4\x8b\x59\x50\xd9\x73\x65\x6b\x88\x22\x2f\xb0\x14\xc0\x1e\xb9\x43\x80\xdb\xec\x2e\x6b\xc0\xb8\xd3\xdb\x3d\x60\xda\x60\x4e\xf2\x9f\xb5\xa5\x83\xa2\x21\xfa\x04\x0f\xf7\x4d\x69\xe3\x0d\xae\x16\x6b\xba\xdd\x39\x21\xaf\xe2\xf0\xeb\x9c\xb6\x58\xdf\xb6\xc7\xc4\x2d\xbc\xfe\x92\xd8\xc1\x52\xb9\x97\x8a\xc5\x3c\x03\x00\x00\xff\xff\xc0\xbe\x81\x75\x07\x01\x00\x00"),
		},
		"/objectstore_editor_role.yaml": &vfsgen۰CompressedFileInfo{
			name:             "objectstore_editor_role.yaml",
			modTime:          time.Time{},
			uncompressedSize: 417,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x8e\xb1\x4e\x03\x31\x0c\x40\xf7\x7c\x85\x25\xe6\x4b\xc5\x86\xb2\x32\xb0\x33\xb0\xfb\x12\xab\x35\xcd\xc5\x91\xed\x14\xa9\x5f\x8f\x72\xb4\x52\x05\x2b\x53\x9e\x92\xe7\xf8\x3d\x41\x27\xdd\xd8\x8c\xa5\x19\xb8\x40\x11\xa0\xc2\x0e\xb2\x7e\x52\x76\x73\x51\xb2\x18\xb0\xf3\x07\xe9\x94\x12\xe8\x8a\x39\xe2\xf0\x93\x28\x5f\xd1\x59\x5a\x3c\xbf\x58\x64\x39\x5c\x9e\xc3\x99\x5b\x49\xf0\x5a\x87\x39\xe9\xbb\x54\x0a\x1b\x39\x16\x74\x4c\x01\xa0\xe1\x46\xe9\xf1\xeb\x65\xee\x12\x5d\x74\x9a\x3a\x2a\x59\x0a\x0b\x60\xe7\x37\x95\xd1\x6d\xce\x2c\xb0\x49\x9b\x12\xb7\x63\x5c\xb1\x5d\x91\x73\x95\x51\x22\x4b\x00\x50\x32\x19\x9a\xe9\xa6\x3e\x56\x07\x80\x0b\xe9\x7a\x7b\xc9\x4a\xe8\xb4\x63\xa1\x4a\x37\x3c\x92\xef\x67\x65\xfb\x81\x8e\x9e\x4f\x3b\x8d\x5e\xee\x03\x5f\xfb\xe5\x7f\x75\x1d\xcc\xd1\xc7\xaf\xbc\x7b\xc8\x9f\xfd\xdf\x01\x00\x00\xff\xff\x16\x16\xea\x35\xa1\x01\x00\x00"),
		},
		"/objectstore_viewer_role.yaml": &vfsgen۰CompressedFileInfo{
			name:             "objectstore_viewer_role.yaml",
			modTime:          time.Time{},
			uncompressedSize: 355,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x8e\xa1\x6e\x03\x31\x0c\x86\x79\x9e\xc2\xd2\xf0\xa5\x1a\x9b\x42\x07\xc6\x07\xc6\x7d\x89\xd5\x7a\x4d\xe2\xc8\x76\xae\x52\x9f\x7e\xba\xb5\xe0\xc0\xe0\x90\x2d\xf9\xf3\xff\xfd\x2f\x30\x48\x1b\x9b\xb1\x74\x03\x17\x28\x02\x1b\xd3\x8d\x14\x64\xfd\xa6\xec\xe6\xa2\x64\x31\xe0\xe0\x2f\xd2\x1d\x4b\xa0\x2b\xe6\x88\xd3\x2f\xa2\x7c\x47\x67\xe9\xf1\xfa\x66\x91\xe5\xb4\xbd\x86\x2b\xf7\x92\xe0\xbd\x4e\x73\xd2\x4f\xa9\x14\x1a\x39\x16\x74\x4c\x01\xa0\x63\xa3\x74\x8c\x5e\x1e\xb6\x45\x77\x52\x67\x25\x4b\x61\x01\x1c\xfc\xa1\x32\x87\xed\x3f\x0b\x34\xe9\xec\xa2\xdc\xcf\x71\xc5\x7e\x47\xce\x55\x66\x89\x2c\x01\x40\xc9\x64\x6a\xa6\x27\x7a\x6c\x1d\x00\x36\xd2\xf5\x79\x39\x93\xff\xce\xca\xf6\x58\x6e\xe8\xf9\xf2\x6f\xb2\x93\x39\xfa\xfc\xc3\xf9\x13\x00\x00\xff\xff\x9c\xb7\x52\xe3\x63\x01\x00\x00"),
		},
		"/role.yaml": &vfsgen۰CompressedFileInfo{
			name:             "role.yaml",
			modTime:          time.Time{},
			uncompressedSize: 2449,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x95\xcd\x8e\xd3\x40\x0c\x80\xef\xf3\x14\xa3\xbd\xa7\x2b\x6e\x28\x57\x0e\xdc\x11\xe2\xee\x4e\xdc\xd4\x74\xc6\x1e\xd9\x9e\x2c\xec\xd3\xa3\xa4\x2d\x88\xfe\x48\x08\x96\xa8\xa7\x58\x4e\x1c\x7f\xf3\xd9\x51\x42\xd7\x75\x01\x2a\x7d\x41\x35\x12\xee\xa3\x6e\x21\x6d\xa0\xf9\x5e\x94\x5e\xc1\x49\x78\x73\x78\x6f\x1b\x92\xe7\xe9\x5d\x38\x10\x0f\x7d\xfc\x90\x9b\x39\xea\x27\xc9\x18\x0a\x3a\x0c\xe0\xd0\x87\x18\x93\xe2\x52\xf0\x99\x0a\x9a\x43\xa9\x7d\xe4\x96\x73\x88\x91\xa1\x60\x1f\x0b\x30\x8c\xa8\x9d\xce\x85\xda\x32\x5a\x1f\xba\x08\x95\x3e\xaa\xb4\x6a\xf3\x2b\xba\xf8\xf4\x14\x62\x54\x34\x69\x9a\xf0\x94\x4b\xc2\x3b\x1a\x0b\x54\x0b\x31\x4e\xa8\xdb\x73\x7e\x6e\x88\x4b\x38\x60\xc6\x53\x38\xa2\x2f\xd7\x4c\x76\x0c\x2a\x78\xda\x2f\x51\xab\xc3\xb9\xe0\x65\x49\xfe\x51\xfb\x3a\xbb\x31\x47\xf6\x49\x72\x2b\x98\x32\x50\xb1\xe5\x96\xa1\x4e\x94\xf0\xff\x73\x41\x25\xfc\xe6\xc8\xf3\x90\xec\x34\x91\x1b\xa2\x9a\xb9\x94\x73\x72\xc0\x1d\x31\xcd\x13\xb9\xe0\xbb\x22\xb9\xdb\x74\x31\x7e\xd1\xc4\x1c\x1c\x77\x2d\x1b\xfa\x1a\xe7\xae\x47\xd3\xbf\x4e\x7f\x4d\x34\x60\xcd\xf2\xbd\x20\xaf\x00\xf4\x1b\x47\x17\x19\xfd\x45\xf4\x40\x3c\xde\x1d\x0a\xf1\xa8\x68\xb6\xc2\x92\x14\x61\x72\xd1\x19\x66\x0b\xfc\x0a\x94\xb2\xb4\xe1\x26\x94\x6c\xbf\x62\x72\x73\xd1\xc7\xe5\x7a\x9e\x37\xad\xdd\xd9\xdd\xab\xfe\xff\xd2\x55\x31\x21\x4d\xa8\x0f\xa5\xe2\x27\xd4\x7a\x1e\x16\xef\xc8\x43\x15\x5a\xe3\x5b\xfa\x6b\xb2\xf5\x8c\xf8\x1e\x58\x1e\xca\xc4\x91\x68\x6d\x03\x0f\xb9\x14\x17\x68\x6f\xe7\x24\x89\xa2\xd8\x26\x49\xb9\xb1\x89\xc7\x5f\xfe\xe9\xe9\x37\xf7\xf1\x23\x00\x00\xff\xff\x0b\xb3\xe9\x93\x91\x09\x00\x00"),
		},
		"/role_binding.yaml": &vfsgen۰CompressedFileInfo{
			name:             "role_binding.yaml",
			modTime:          time.Time{},
			uncompressedSize: 261,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x8d\xb1\x4e\x03\x31\x0c\x86\xf7\x3c\x85\x5f\xa0\x87\xd8\x50\x36\x60\x60\x2f\x12\xbb\x2f\x71\x8b\x69\x62\x47\xb6\x53\x09\x9e\x1e\x9d\x74\xb0\x80\xba\xd9\xd2\xf7\xfd\x1f\x0e\x7e\x23\x73\x56\xc9\x60\x2b\x96\x05\x67\xbc\xab\xf1\x17\x06\xab\x2c\x97\x07\x5f\x58\xef\xae\xf7\xe9\xc2\x52\x33\x3c\xb7\xe9\x41\x76\xd4\x46\x4f\x2c\x95\xe5\x9c\x3a\x05\x56\x0c\xcc\x09\x40\xb0\x53\x86\x8e\x82\x67\xb2\x83\x69\xa3\x75\xa7\xb6\xfb\x48\xa7\x0d\xc2\xc1\x2f\xa6\x73\xdc\x08\x26\x80\x3f\xbd\x7f\xe7\x93\xcf\xf5\x83\x4a\x78\x4e\x87\x5d\x79\x25\xbb\x72\xa1\xc7\x52\x74\x4a\xfc\x5a\x95\x4e\x38\xdb\xcf\xef\x03\x0b\x65\xf0\x4f\x0f\xea\xe9\x3b\x00\x00\xff\xff\x3d\x57\xd3\xb6\x05\x01\x00\x00"),
		},
		"/storeendpoint_editor_role.yaml": &vfsgen۰CompressedFileInfo{
			name:             "storeendpoint_editor_role.yaml",
			modTime:          time.Time{},
			uncompressedSize: 425,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x8e\x31\x6f\x33\x21\x0c\x40\x77\x7e\x85\xa5\x6f\x3e\xa2\x6f\xab\x58\x3b\x74\xef\xd0\xdd\x07\x56\x62\x85\xc3\xc8\x36\xa9\x94\x5f\x5f\x71\xbd\x0c\x69\xd7\x4e\x3c\xc1\x33\x7e\xff\xa0\x93\x6e\x6c\xc6\xd2\x0c\x5c\xa0\x08\x50\x61\x07\x73\x51\xa2\x56\xba\x70\x73\x8b\x01\x3b\x7f\x90\x4e\x2d\x81\xae\x98\x23\x0e\xbf\x88\xf2\x1d\x9d\xa5\xc5\xeb\x8b\x45\x96\xd3\xed\x7f\xb8\x72\x2b\x09\x5e\xeb\x30\x27\x7d\x97\x4a\x61\x23\xc7\x82\x8e\x29\x00\x34\xdc\x28\x3d\x7f\xbe\xcc\x7d\xa2\x8b\x4e\x57\x47\x25\x4b\x61\x01\xec\xfc\xa6\x32\xba\xcd\xa9\x05\x36\x69\x53\xe2\x76\x8e\x2b\xb6\x3b\x72\xae\x32\x4a\x64\x09\x00\x4a\x26\x43\x33\x1d\xea\x73\x79\x00\xb8\x91\xae\xc7\x5b\x56\x42\xa7\x1d\x0b\x55\x3a\xf0\x4c\xbe\x9f\x95\xed\x1b\x3a\x7a\xbe\xec\x34\x7a\x79\x0c\x7c\xee\x97\x7f\x57\x76\x32\x47\x1f\x3f\x02\x1f\x29\xbf\x0a\xbe\x02\x00\x00\xff\xff\x6e\x15\x06\xe0\xa9\x01\x00\x00"),
		},
		"/storeendpoint_viewer_role.yaml": &vfsgen۰CompressedFileInfo{
			name:             "storeendpoint_viewer_role.yaml",
			modTime:          time.Time{},
			uncompressedSize: 363,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x8e\xb1\x6a\xc3\x40\x0c\x86\xf7\x7b\x0a\x41\x67\x5f\xe8\x56\x6e\xed\xd0\xbd\x43\x77\xd9\x27\x12\x91\xb3\x74\x48\x3a\x07\xf2\xf4\xc5\x4d\x16\x43\xc7\x4c\x12\xe8\xd3\xff\xfd\x6f\xd0\xc9\x56\x76\x67\x15\x87\x50\xa8\x0a\x1b\xd3\x8d\x0c\x3c\xd4\x88\xa4\x76\x65\x09\xcf\x09\x3b\xff\x90\xed\x60\x01\x9b\x71\xc9\x38\xe2\xa2\xc6\x77\x0c\x56\xc9\xd7\x0f\xcf\xac\xa7\xed\x3d\x5d\x59\x6a\x81\xcf\x36\x3c\xc8\xbe\xb5\x51\x5a\x29\xb0\x62\x60\x49\x00\x82\x2b\x95\x63\xf8\xf4\x30\x4e\xb6\xb3\x36\x1a\x79\x49\x13\x60\xe7\x2f\xd3\xd1\x7d\xff\x9a\x60\x55\xe1\x50\x63\x39\xe7\x19\xe5\x8e\xbc\x34\x1d\x35\xb3\x26\x00\x23\xd7\x61\x0b\x3d\xd1\x63\xf3\x04\xb0\x91\xcd\xcf\xdb\x99\xe2\x6f\x36\xf6\xc7\x72\xc3\x58\x2e\x2f\xd4\x9d\x3c\x30\xc6\x3f\xd6\xdf\x00\x00\x00\xff\xff\xe6\x13\x59\x6f\x6b\x01\x00\x00"),
		},
		"/thanos_editor_role.yaml": &vfsgen۰CompressedFileInfo{
			name:             "thanos_editor_role.yaml",
			modTime:          time.Time{},
			uncompressedSize: 394,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x8e\xb1\x4e\x03\x31\x0c\x86\xf7\x3c\x85\x25\xe6\x4b\xc5\x86\xb2\x32\xb0\x33\xb0\xbb\x89\xd5\xb3\x9a\x8b\x23\xdb\x29\x52\x9f\x1e\xe5\x7a\x5d\x60\x64\xca\x27\xe7\xb3\xff\xff\x05\x3a\xe9\xc6\x66\x2c\xcd\xc0\x05\x8a\x00\x15\x76\xf0\x15\x9b\x58\x0c\xd8\xf9\x8b\x74\x7e\x27\xd0\x33\xe6\x88\xc3\x57\x51\xbe\xa3\xb3\xb4\x78\x7d\xb3\xc8\x72\xba\xbd\x86\x2b\xb7\x92\xe0\xbd\x0e\x73\xd2\x4f\xa9\x14\x36\x72\x2c\xe8\x98\x02\x40\xc3\x8d\xd2\x71\x74\x99\x01\xa2\x8b\x4e\x49\x47\x25\x4b\x61\x01\xec\xfc\xa1\x32\xba\x4d\x7d\x81\x4d\xda\x94\xb8\x5d\xe2\x19\xdb\x1d\x39\x57\x19\x25\xb2\x04\x00\x25\x93\xa1\x99\x0e\xf5\x71\x35\x00\xdc\x48\xcf\xc7\x2c\x2b\xa1\xd3\x8e\x85\x2a\x1d\x78\x21\xdf\xdf\xca\xf6\x80\x8e\x9e\xd7\x9d\x46\x2f\xcf\x85\xef\x7d\xf8\xff\x46\x27\x73\xf4\xf1\xab\xd8\xb3\xc2\x9f\xe4\x9f\x00\x00\x00\xff\xff\xe0\xf8\xd9\x2c\x8a\x01\x00\x00"),
		},
		"/thanos_viewer_role.yaml": &vfsgen۰CompressedFileInfo{
			name:             "thanos_viewer_role.yaml",
			modTime:          time.Time{},
			uncompressedSize: 332,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x8d\x31\x6b\xc4\x30\x0c\x46\x77\xff\x0a\x41\xe7\xf8\xe8\x56\xbc\x76\xe8\xde\xa1\xbb\x92\x88\x8b\x38\x47\x32\x92\x9c\x83\xfc\xfa\x92\x26\x63\xc7\x9b\x24\x3e\x1e\xef\xbd\x41\x23\x5b\xd9\x9d\x55\x1c\x42\x61\x56\xd8\x98\x9e\x64\x10\x0b\x8a\x7a\x4e\xd8\xf8\x87\xec\x00\x0a\xd8\x88\x53\xc6\x1e\x8b\x1a\xef\x18\xac\x92\x1f\x1f\x9e\x59\x6f\xdb\x7b\x7a\xb0\xcc\x05\x3e\x6b\xf7\x20\xfb\xd6\x4a\x69\xa5\xc0\x19\x03\x4b\x02\x10\x5c\xa9\x5c\xd2\xe1\x4c\x0c\x76\x40\xd6\x2b\x79\x49\x03\x60\xe3\x2f\xd3\xde\xfc\xc0\x07\x58\x55\x38\xd4\x58\xee\x79\x44\xd9\x91\xa7\xaa\x7d\xce\xac\x09\xc0\xc8\xb5\xdb\x44\x17\x7a\x5a\x13\xc0\x46\x36\x5e\xdb\x9d\xe2\xef\x56\xf6\xf3\x79\x62\x4c\xcb\x0b\x32\x37\x0f\x8c\xfe\x4f\xed\x37\x00\x00\xff\xff\xb4\x8f\x3e\x8e\x4c\x01\x00\x00"),
		},
	}
	fs["/"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/auth_proxy_role.yaml"].(os.FileInfo),
		fs["/auth_proxy_role_binding.yaml"].(os.FileInfo),
		fs["/auth_proxy_service.yaml"].(os.FileInfo),
		fs["/kustomization.yaml"].(os.FileInfo),
		fs["/leader_election_role.yaml"].(os.FileInfo),
		fs["/leader_election_role_binding.yaml"].(os.FileInfo),
		fs["/objectstore_editor_role.yaml"].(os.FileInfo),
		fs["/objectstore_viewer_role.yaml"].(os.FileInfo),
		fs["/role.yaml"].(os.FileInfo),
		fs["/role_binding.yaml"].(os.FileInfo),
		fs["/storeendpoint_editor_role.yaml"].(os.FileInfo),
		fs["/storeendpoint_viewer_role.yaml"].(os.FileInfo),
		fs["/thanos_editor_role.yaml"].(os.FileInfo),
		fs["/thanos_viewer_role.yaml"].(os.FileInfo),
	}

	return fs
}()

type vfsgen۰FS map[string]interface{}

func (fs vfsgen۰FS) Open(path string) (http.File, error) {
	path = pathpkg.Clean("/" + path)
	f, ok := fs[path]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}
	}

	switch f := f.(type) {
	case *vfsgen۰CompressedFileInfo:
		gr, err := gzip.NewReader(bytes.NewReader(f.compressedContent))
		if err != nil {
			// This should never happen because we generate the gzip bytes such that they are always valid.
			panic("unexpected error reading own gzip compressed bytes: " + err.Error())
		}
		return &vfsgen۰CompressedFile{
			vfsgen۰CompressedFileInfo: f,
			gr:                        gr,
		}, nil
	case *vfsgen۰DirInfo:
		return &vfsgen۰Dir{
			vfsgen۰DirInfo: f,
		}, nil
	default:
		// This should never happen because we generate only the above types.
		panic(fmt.Sprintf("unexpected type %T", f))
	}
}

// vfsgen۰CompressedFileInfo is a static definition of a gzip compressed file.
type vfsgen۰CompressedFileInfo struct {
	name              string
	modTime           time.Time
	compressedContent []byte
	uncompressedSize  int64
}

func (f *vfsgen۰CompressedFileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *vfsgen۰CompressedFileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *vfsgen۰CompressedFileInfo) GzipBytes() []byte {
	return f.compressedContent
}

func (f *vfsgen۰CompressedFileInfo) Name() string       { return f.name }
func (f *vfsgen۰CompressedFileInfo) Size() int64        { return f.uncompressedSize }
func (f *vfsgen۰CompressedFileInfo) Mode() os.FileMode  { return 0444 }
func (f *vfsgen۰CompressedFileInfo) ModTime() time.Time { return f.modTime }
func (f *vfsgen۰CompressedFileInfo) IsDir() bool        { return false }
func (f *vfsgen۰CompressedFileInfo) Sys() interface{}   { return nil }

// vfsgen۰CompressedFile is an opened compressedFile instance.
type vfsgen۰CompressedFile struct {
	*vfsgen۰CompressedFileInfo
	gr      *gzip.Reader
	grPos   int64 // Actual gr uncompressed position.
	seekPos int64 // Seek uncompressed position.
}

func (f *vfsgen۰CompressedFile) Read(p []byte) (n int, err error) {
	if f.grPos > f.seekPos {
		// Rewind to beginning.
		err = f.gr.Reset(bytes.NewReader(f.compressedContent))
		if err != nil {
			return 0, err
		}
		f.grPos = 0
	}
	if f.grPos < f.seekPos {
		// Fast-forward.
		_, err = io.CopyN(ioutil.Discard, f.gr, f.seekPos-f.grPos)
		if err != nil {
			return 0, err
		}
		f.grPos = f.seekPos
	}
	n, err = f.gr.Read(p)
	f.grPos += int64(n)
	f.seekPos = f.grPos
	return n, err
}
func (f *vfsgen۰CompressedFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.seekPos = 0 + offset
	case io.SeekCurrent:
		f.seekPos += offset
	case io.SeekEnd:
		f.seekPos = f.uncompressedSize + offset
	default:
		panic(fmt.Errorf("invalid whence value: %v", whence))
	}
	return f.seekPos, nil
}
func (f *vfsgen۰CompressedFile) Close() error {
	return f.gr.Close()
}

// vfsgen۰DirInfo is a static definition of a directory.
type vfsgen۰DirInfo struct {
	name    string
	modTime time.Time
	entries []os.FileInfo
}

func (d *vfsgen۰DirInfo) Read([]byte) (int, error) {
	return 0, fmt.Errorf("cannot Read from directory %s", d.name)
}
func (d *vfsgen۰DirInfo) Close() error               { return nil }
func (d *vfsgen۰DirInfo) Stat() (os.FileInfo, error) { return d, nil }

func (d *vfsgen۰DirInfo) Name() string       { return d.name }
func (d *vfsgen۰DirInfo) Size() int64        { return 0 }
func (d *vfsgen۰DirInfo) Mode() os.FileMode  { return 0755 | os.ModeDir }
func (d *vfsgen۰DirInfo) ModTime() time.Time { return d.modTime }
func (d *vfsgen۰DirInfo) IsDir() bool        { return true }
func (d *vfsgen۰DirInfo) Sys() interface{}   { return nil }

// vfsgen۰Dir is an opened dir instance.
type vfsgen۰Dir struct {
	*vfsgen۰DirInfo
	pos int // Position within entries for Seek and Readdir.
}

func (d *vfsgen۰Dir) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		d.pos = 0
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported Seek in directory %s", d.name)
}

func (d *vfsgen۰Dir) Readdir(count int) ([]os.FileInfo, error) {
	if d.pos >= len(d.entries) && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(d.entries)-d.pos {
		count = len(d.entries) - d.pos
	}
	e := d.entries[d.pos : d.pos+count]
	d.pos += count
	return e, nil
}
