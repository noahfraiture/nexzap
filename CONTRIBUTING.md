
# Contributing to Nexzap

Thanks for wanting to contribute to Nexzap, my open-source platform for learning a new programming language every week! Below, I outline how you can contribute code and tutorials. Keep it simple and let’s make something awesome together.

## Code Contributions

1. **Fork the Repository**: Fork the Nexzap repo on GitHub and clone your fork.
2. **Create a Branch**: Use a descriptive branch name (e.g., `feature/add-auth` or `fix/bug-login`).
3. **Follow Coding Standards**:
   - Stick to Go, HTMX, Go Templ, TailwindCSS, DaisyUI, and Alpine.js.
   - Keep code clean and documented.
4. **PR Name and Commit**:
   - Use clear commit messages with a type prefix like: `feat:`, `fix:`, `tuto:`, or `refactor:`, followed by a short description (e.g., `feat: add user auth`, `fix: login error`, `tuto: add Python tutorial`).
   - Avoid vague or funny commit messages in the main branch.
   - PR titles should also be descriptive and follow a similar format.
5. **Test Your Changes**: Ensure your changes work locally and do not introduce errors.
   - **Set Up PostgreSQL**:
     - Run a PostgreSQL instance using `docker compose -f compose-local.yml up -d`.
     - Verify the instance is running with `docker ps`.
   - **Install Frontend Dependencies**:
     - Install TailwindCSS and DaisyUI by running `npm install` in the project root.
   - **Install Backend Tools**:
     - Install `go templ` and `sqlc` using your preferred method:
       - **Standard Installation**: Follow the official guides:
         - [Templ Installation Guide](https://templ.guide/quick-start/installation)
         - [Sqlc Installation Guide](https://docs.sqlc.dev/en/latest/overview/install.html)
       - **NixOS Users**: Run `nix-shell ./flake.nix` to enter a shell with `go templ` and `sqlc` preconfigured.
     - **Note**: 
       - `sqlc` generates SQL boilerplate code and is only required if modifying SQL queries.
       - `go templ` compiles `.templ` files for frontend rendering and is only needed if working on frontend templates.
   - **Run Build and Test Commands**:
     - Use `go-task` to execute all necessary tasks with a single command: `task default`.
     - Alternatively, run individual commands manually as defined in `./Taskfile.yml`:
       - Compile TailwindCSS: `tailwindcss -i ./static/css/input.css -o ./static/css/style.css --minify`
       - Generate templates: `templ generate`
       - Generate SQL code: `sqlc generate`
       - Run the application: `go run ./cmd/nexzap`
     - Verify the application runs locally by accessing it in your browser (default: `http://localhost:8080`).
6. **Submit a Pull Request (PR)**:
   - Include a clear description of your changes.
   - Reference related issues.
   - Ensure all checks pass before requesting a review.

## Tutorial Contributions

I’m excited to accept community tutorials for Nexzap! Each tutorial introduces a programming language through its core concepts, written in a light, accessible, and fun way.

### How to Write a Tutorial

Check out the `./tutorials` directory to see how tutorials are structured. Here’s the setup:

- **Tutorial Directory**: Each tutorial gets its own folder. The name doesn’t affect the code, but I’d prefer names in the order they’re written (e.g., `0_go`, `1_python`).
- **Directory Contents**:
   1. **Sheet Folders** (e.g., `1_overview`): These contain the tutorial’s sheets. Name them with a number prefix starting at `1` (e.g., `1_overview`, `2_functions`), with no gaps in numbering.
   2. **`meta.toml`**: Includes metadata like the tutorial’s title, CodeMirror CSS library for code highlighting, version (set to `1` for new tutorials), and the unlock date.
      ```toml
      title = "Go"
      codeEditor = "go"
      version = 1
      unlock = "2025-04-25"
      ```
   You can find CodeMirror mode for language [here](https://cdnjs.com/libraries/codemirror/5.65.18).
   It's most likely that I will change the unlock date do fit my schedule. However feel free to discuss.

   3. **`docker/`**: Contains a `Dockerfile` to build the base image for testing code.

- **Sheet Folder Structure** (e.g., `1_overview`):
   - **`correction/`**: Holds files copied to the container for testing.
   - **Placeholder File**: A file (e.g., `main.go`) with starter code for the exercise. User input replaces its content in `correction/` during testing.
   - **`exercise.md`**: Instructions for the exercise.
   - **`guide.md`**: The guide content shown in the left panel, introducing the language or concept.
   - **`meta.toml`**: Specifies the Docker image name (I build these, myself for now. The name is flexible but you can put the same name as the name of the tutorial directory), the test command, and the placeholder file name.
      ```toml
      image = "gotest"
      command = "go test"
      submission = "main.go"
      ```

### Tutorial Submission Process

1. **Create a Tutorial**:
   - Write it in Markdown, following the structure above.
   - Include a `meta.toml` with an `unlock` date (e.g., `2025-06-01`).
2. **Submit a Pull Request**:
   - Place your tutorial in the `tutorials/` directory.
   - Describe the language and concepts in the PR description.
   - Use the `tuto:` prefix for commits (e.g., `tuto: add Python tutorial`).
3. **Review Process**:
   - I’ll review tutorials for clarity, accuracy, and alignment with Nexzap’s fun, accessible vibe.
   - I reserve the right to:
      - Change the `unlock` date at any time, even after merging a PR.
      - Decide not to publish a tutorial on the website, even if the PR is merged or the `unlock` date has passed.
   - Accepted tutorials are manually deployed to the server and added to the database. They’re only visible to users after the `unlock` date and once I’ve deployed them.

## Important Notes

I explicitly reserve the right to control tutorial publication, including modifying `unlock` dates or choosing not to publish tutorials, regardless of PR status or date. If I accept your PR you will most likely have your tutorial on the website at some point, but for question of organization, there's no 100% garantees.

If you have any question, feel free to ask.

Thanks for helping make Nexzap a fun way to learn programming languages!
