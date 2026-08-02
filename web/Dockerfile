# Use an official Node runtime as a parent image
FROM node:20

# Set the working directory in the container
WORKDIR /usr/src/app

RUN mkdir cert

# Install pnpm first so this layer is not invalidated by app changes.
# Pinned to 9.x to match the lockfileVersion 9.0 in pnpm-lock.yaml.
RUN npm install -g pnpm@9

# Copy only the manifest and lockfile so dependency install is cached
# and is not invalidated by unrelated source changes.
COPY package.json pnpm-lock.yaml ./

# Install app dependencies straight from the lockfile (no registry re-resolution)
RUN pnpm install --frozen-lockfile

# Copy the rest of the application code to the working directory
COPY . .

# Build the app
RUN pnpm run build

# Expose the port that the app will run on
EXPOSE 3000

RUN chmod +x ./run.sh

# Define the command to run the app
CMD ["./run.sh"]