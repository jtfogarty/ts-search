# README

## About

This project was created using a custom creation script (`creationScript.sh`) that sets up a Wails project with Svelte, TailwindCSS, and shadcn-svelte components. The script automates the process of initializing a Wails project, configuring Svelte, and adding specified UI components.

## Project Creation

The project was created using the following command:

```bash
./creationScript.sh bun app01 'button,card,input'
```

### Creation Script

<details>
<summary>Click to view the creation script</summary>

```bash
#!/bin/bash

# Check if all required arguments are provided
if [ "$#" -lt 3 ]; then
    echo "Usage: $0 <package_manager> <project_name> <components> [brand]"
    echo "Example: $0 bun my_project 'button,card,input' mybrand"
    exit 1
fi

manager=$1
project=$2
components=$3
brand=$4

# Initialize Wails project
wails init -n $project -t svelte
cd $project

# Update wails.json
sed -i '' "s|npm|$manager|g" wails.json
sed -i '' 's|"auto",|"auto",\n  "wailsjsdir": "./frontend/src/lib",|' wails.json

# Update main.go
sed -i '' "s|all:frontend/dist|all:frontend/build|" main.go

# Handle branding if specified
if [[ -n $brand ]]; then
    mv frontend/src/App.svelte +page.svelte
    sed -i '' "s|'./assets|'\$lib/assets|" +page.svelte
    sed -i '' "s|'../wails|'\$lib/wails|" +page.svelte
    mv frontend/src/assets .
fi

# Remove old frontend and create new Svelte app
rm -r frontend
$manager x create-svelte@latest frontend

# Move files for branding
if [[ -n $brand ]]; then
    mv +page.svelte frontend/src/routes/+page.svelte
    mkdir -p frontend/src/lib
    mv assets frontend/src/lib/
fi

# Install dependencies and configure Svelte
cd frontend
$manager install
$manager uninstall @sveltejs/adapter-auto
$manager add -D @sveltejs/adapter-static

# Add TailwindCSS
$manager x svelte-add@latest tailwindcss
$manager install

# Setup shadcn-svelte
$manager add -D shadcn-svelte
$manager x shadcn-svelte@latest init

# Add specified shadcn components
IFS=',' read -ra ADDR <<< "$components"
for component in "${ADDR[@]}"; do
    $manager x shadcn-svelte@latest add $component
done

# Create +layout.ts
echo "export const prerender = true;" > src/routes/+layout.ts
echo "export const ssr = false;" >> src/routes/+layout.ts

# Update svelte.config.js
cat > svelte.config.js << EOL
import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
    preprocess: vitePreprocess(),
    kit: {
        adapter: adapter({
            pages: 'build',
            assets: 'build',
            fallback: null,
            precompress: false,
            strict: true
        }),
        alias: {
            "\$lib": "./src/lib",
            "@/*": "./src/lib/*"
        }
    }
};

export default config;
EOL

# Return to project root
cd ..

echo "Setup complete. I will now run 'wails dev' to start the development server."

# Start Wails development server
wails dev
```

</details>
</br>
This command initializes a Wails project named "app01" using Bun as the package manager and includes the 'button', 'card', and 'input' components from shadcn-svelte.

During the project creation, the following options were selected:

1. Svelte app template: Skeleton project
2. Type checking: Yes, using TypeScript syntax
3. Additional options:
   - Add ESLint for code linting
   - Add Prettier for code formatting
   - Try the Svelte 5 preview (unstable!)

These selections ensure a robust development environment with type safety, code quality tools, and the latest Svelte features.


## TailwindCSS Integration

TailwindCSS was added to the project using the svelte-add tool. This integration provides a utility-first CSS framework that can be composed to build any design, directly in your markup.

## shadcn-svelte Components

The project incorporates shadcn-svelte, a collection of re-usable components built using Radix and Tailwind CSS. The following components were added during project creation:

- Button
- Card
- Input

These components provide a solid foundation for building a user interface with consistent design and functionality.

### shadcn-svelte Configuration

During the setup process, shadcn-svelte was initialized with default options. This includes:

- Setting up the necessary directory structure in `src/lib/components/ui`
- Configuring the theme in `src/app.postcss`
- Adding utility functions in `src/lib/utils.ts`

The components are ready to use and can be imported from their respective locations in the `src/lib/components/ui` directory.


## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.
