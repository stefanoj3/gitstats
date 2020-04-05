// Package usecase provides application use cases.
//
// For the Dependency Rule it is not allowed to import platform.
//
// Something to clean up in the future: the current usecase for getting statistics is tightly coupled with the github
// library and thus coupled to github data structures.
// We should consider to either use generic data structures so we abstract away from the platform and the library,
// or we should make the usecase github specific.
package usecase
