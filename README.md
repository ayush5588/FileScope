# FileScope

FileScope is a web application that helps you find all GitHub Pull Requests (PRs), both open and closed, that have modified a specific file in a repository. Simply provide the full GitHub file URL, and FileScope will return a list of PRs that have changed that file.

## Features
- **Find PRs by File:** Enter a GitHub file URL to see all PRs that have modified it.
- **Web UI:** Simple web interface for quick searching.
- **API:** Programmatic access to the core functionality.
- **Supports Open & Closed PRs:** See all relevant PRs, not just open ones.

## Getting Started

### Prerequisites
- Go 1.20+
- A GitHub personal access token (for higher rate limits)

### Setup
1. **Clone the repository:**
   ```sh
   git clone https://github.com/ayush5588/FileScope.git
   cd FileScope
   ```
2. **Set up your GitHub token:**
   - Create a file named `token.env` in the root directory.
   - Add your GitHub token as an environment variable:
     ```env
     GH_TOKEN=your_github_token_here
     ```
   - (Optional) For higher rate limits, you can add multiple tokens as `GITHUB_TOKEN_1`, `GITHUB_TOKEN_2`, etc., and set `TOKEN_COUNT` accordingly.

3. **Run the application:**
   ```sh
   go run main.go
   ```
   The server will start on [http://localhost:8080](http://localhost:8080).

### Docker
You can also run FileScope using Docker:
```sh
docker build -t filescope .
docker run -p 8080:8080 --env-file token.env filescope
```

## Usage
- Open [http://localhost:8080](http://localhost:8080) in your browser.
- Enter the full GitHub file URL (e.g., `https://github.com/owner/repo/blob/branch/path/to/file.go`).
- Submit to see a table of PRs that have modified the file.

## API Endpoints
- `GET /healthz` — Health check endpoint.
- `POST /getPR` — Accepts a form field `filePath` (the GitHub file URL) and returns a JSON list of PRs that have modified the file.

## Example Request
```sh
curl -X POST -F 'filePath=https://github.com/owner/repo/blob/branch/path/to/file.go' http://localhost:8080/getPR
```

## Output Format
Each PR in the response includes:
- PR number (with link)
- Creation date
- Creator
- Title
- Branch

## Error Handling
If the input URL is invalid or the file cannot be found, a helpful error message will be returned.

## License
MIT 
