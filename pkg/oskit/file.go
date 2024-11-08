package oskit

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func InstallDirectory(sourceDir, destDir string) {
	// Create destination base directory
	err := os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create destination directory: %v", err)
	}

	// Walk through the source directory
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing %s: %v", path, err)
		}

		// Define the corresponding destination path
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("failed to calculate relative path: %v", err)
		}
		destPath := filepath.Join(destDir, relPath)

		if info.IsDir() {
			// Create the directory in the destination
			err := os.MkdirAll(destPath, info.Mode())
			if err != nil {
				return fmt.Errorf("failed to create directory %s: %v", destPath, err)
			}
			return nil
		}

		// Handle file based on its permissions
		if info.Mode()&0o111 != 0 {
			// File has execute permission, copy with permissions
			err := copyFileWithPermissions(path, destPath, info.Mode())
			if err != nil {
				return fmt.Errorf("error copying executable file %s: %v", path, err)
			}
		} else {
			// File without execute permission, simple copy
			err := copyFile(path, destPath)
			if err != nil {
				return fmt.Errorf("error copying non-executable file %s: %v", path, err)
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Error processing files: %v", err)
	}

	fmt.Println("Files copied successfully.")
}

// Copy file with permissions (for executable files)
func copyFileWithPermissions(src, dst string, perm os.FileMode) error {
	// Ensure the destination directory exists
	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create destination directory for %s: %v", dst, err)
	}

	// Copy the file content
	if err := copyFile(src, dst); err != nil {
		return err
	}

	// Set the permissions for executable files
	if err := os.Chmod(dst, perm); err != nil {
		return fmt.Errorf("failed to set permissions on %s: %v", dst, err)
	}
	return nil
}

// Simple file copy function
func copyFile(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %v", src, err)
	}
	defer srcFile.Close()

	// Ensure the destination directory exists
	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create destination directory for %s: %v", dst, err)
	}

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %v", dst, err)
	}
	defer dstFile.Close()

	// Copy contents from source to destination
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("error copying contents to %s: %v", dst, err)
	}

	return nil
}
