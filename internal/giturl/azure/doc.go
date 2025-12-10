// SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON
// SPDX-License-Identifier: Apache-2.0

// Package azure provides URL parsing and raw content URL generation for Azure DevOps.
//
// # Azure DevOps URL Structure
//
// Azure DevOps uses a more complex URL structure compared to other Git hosting providers.
//
// ## Hosts
//
// Azure DevOps can be hosted on:
//   - dev.azure.com (primary)
//   - prod.azure.com
//   - *.azure.com (custom instances)
//   - ssh.dev.azure.com (SSH access)
//
// ## Path Format
//
// HTTP/HTTPS URLs follow the pattern:
//
//	https://dev.azure.com/{owner}/{project}/_git/{repo}
//
// Key differences from other providers:
//   - Three-part hierarchy: owner (organization), project, and repository
//   - Uses _git as a separator (similar to GitLab's /-/)
//   - Project is a separate entity, not part of the repository path
//
// SSH URLs follow the pattern:
//
//	ssh://git@ssh.dev.azure.com/v3/{owner}/{project}/{repo}
//
// Or shortened:
//
//	git@ssh.dev.azure.com:v3/{owner}/{project}/{repo}
//
// ## Query Parameters
//
// Azure DevOps uses query parameters instead of path segments for version and file path:
//
//   - path: File or directory path
//     Example: ?path=/src/main.go
//
//   - version: Branch or tag reference with prefixes
//     Branch: ?version=GBbranch-name (GB = Git Branch)
//     Tag: ?version=GTtag-name (GT = Git Tag)
//     Examples:
//     ?version=GBmain
//     ?version=GBdev
//     ?version=GTv1.0.0
//
//   - _a: Action parameter (e.g., contents, history)
//     Typically ignored for parsing purposes
//
//   - api-version: API version for REST calls
//     Used when accessing Azure DevOps REST API
//
// ## Example URLs
//
// Repository only:
//
//	https://dev.azure.com/dwertent/ks-testing-public/_git/ks-testing-public
//
// With file path:
//
//	https://dev.azure.com/dwertent/ks-testing-public/_git/ks-testing-public?path=/rules-tests/alert.json
//
// With branch and path:
//
//	https://dev.azure.com/dwertent/ks-testing-public/_git/ks-testing-public?path=/scripts&version=GBdev
//
// With tag and path:
//
//	https://dev.azure.com/dwertent/ks-testing-public/_git/ks-testing-public?path=/scripts&version=GTv1.0.1
//
// Full URL with action:
//
//	https://dev.azure.com/dwertent/ks-testing-public/_git/ks-testing-public?path=/scripts&version=GBdev&_a=contents
//
// # Raw Content URLs
//
// Unlike GitHub, GitLab, Gitea, and Bitbucket which provide simple path-based raw content URLs,
// Azure DevOps uses an API-based approach for accessing raw file content.
//
// ## Azure DevOps Items API
//
// Raw file content is accessed via the Git Items API:
//
//	https://dev.azure.com/{org}/{project}/_apis/git/repositories/{repo}/items
//	  ?path={path}
//	  &version={version}
//	  &api-version=7.0
//	  &download=true
//
// For branch-specific requests:
//
//	https://dev.azure.com/{org}/{project}/_apis/git/repositories/{repo}/items
//	  ?path={path}
//	  &versionDescriptor.version={branch}
//	  &versionDescriptor.versionType=branch
//	  &api-version=7.0
//	  &download=true
//
// For tag-specific requests:
//
//	https://dev.azure.com/{org}/{project}/_apis/git/repositories/{repo}/items
//	  ?path={path}
//	  &versionDescriptor.version={tag}
//	  &versionDescriptor.versionType=tag
//	  &api-version=7.0
//	  &download=true
//
// ## API Version
//
// Azure DevOps REST API versions:
//   - 7.0 (current stable)
//   - 7.1-preview (preview)
//   - 6.0 (older, still supported)
//   - 5.1 (legacy)
//
// The api-version parameter is mandatory for all API calls.
//
// ## Download Parameter
//
// The download=true parameter is crucial:
//   - Without it: Returns JSON metadata about the item
//   - With it: Returns the raw file content directly
//
// # Implementation Challenges
//
// ## Complexity Compared to Other Providers
//
// Azure DevOps raw URLs are more complex because:
//
// 1. API-based instead of path-based
//    - Other providers: https://host/owner/repo/raw/branch/file
//    - Azure: https://host/{org}/{project}/_apis/git/repositories/{repo}/items?...
//
// 2. Query parameter requirements
//    - Multiple query parameters needed (path, version, api-version, download)
//    - Parameters must be properly URL-encoded
//    - Version descriptor needs type specification (branch/tag/commit)
//
// 3. Three-part hierarchy
//    - Organization, project, and repository are all separate entities
//    - Project name may differ from repository name
//    - Must preserve all three parts in the raw URL
//
// 4. Version prefix handling
//    - Browser URLs use GB/GT prefixes (GBmain, GTv1.0.0)
//    - API expects plain version names (main, v1.0.0)
//    - Conversion required between browser format and API format
//
// 5. Authentication
//    - Azure DevOps may require authentication for private repositories
//    - API calls might need PAT (Personal Access Token) in headers
//    - Unlike GitHub raw.githubusercontent.com which works without auth for public repos
//
// ## Recommended Implementation Strategy
//
// When implementing the Azure provider:
//
// 1. Parse Function
//    - Extract owner, project, repo from path (split on /_git/)
//    - Parse query parameters for path and version
//    - Strip GB/GT prefixes from version parameter
//    - Determine version type (branch vs tag vs commit)
//
// 2. Raw Function
//    - Build API URL with /_apis/git/repositories/{repo}/items
//    - Add required query parameters:
//      * path (from locator)
//      * versionDescriptor.version (from locator, without GB/GT)
//      * versionDescriptor.versionType (branch/tag/commit)
//      * api-version=7.0
//      * download=true
//    - Properly URL-encode all parameters
//
// 3. Version Type Detection
//    - If version starts with "GB": branch
//    - If version starts with "GT": tag
//    - If version is SHA-like (40 hex chars): commit
//    - If version is empty: default to main branch
//
// 4. Edge Cases
//    - Empty path: should error (can't fetch repository root as file)
//    - Empty version: default to "main" branch
//    - SSH URLs: convert to HTTPS for API access
//    - Custom Azure instances: may have different API paths
//
// # Testing Considerations
//
// ## Test URLs
//
// Include tests for:
//   - Repository-only URLs
//   - URLs with file paths
//   - URLs with branches (GB prefix)
//   - URLs with tags (GT prefix)
//   - URLs with _a=contents parameter
//   - SSH URLs
//   - Custom Azure DevOps instances
//
// ## Expected Behavior
//
// Valid cases should convert to API URLs with:
//   - Correct api-version parameter
//   - download=true parameter
//   - Properly encoded path parameter
//   - Correct versionDescriptor fields
//
// Invalid cases should reject:
//   - Empty paths (no file to fetch)
//   - Non-Azure hosts
//   - Malformed paths (missing /_git/ separator)
//   - Invalid version prefixes
//
// # References
//
// Azure DevOps REST API documentation:
//   - https://learn.microsoft.com/en-us/rest/api/azure/devops/git/items/get
//   - https://learn.microsoft.com/en-us/rest/api/azure/devops/git/
//
// Example implementation:
//   - https://github.com/armosec/go-git-url/tree/master/azureparser/v1
//
// # TODO
//
// Future implementation tasks:
//   - [ ] Implement Parse function for Azure DevOps URLs
//   - [ ] Implement Raw function using Items API
//   - [ ] Add comprehensive test coverage
//   - [ ] Handle authentication requirements
//   - [ ] Support custom Azure DevOps Server instances
//   - [ ] Document API version compatibility
package azure
