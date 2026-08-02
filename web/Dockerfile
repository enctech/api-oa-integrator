# Use an official Node runtime as a parent image
FROM node:20

# Set the working directory in the container
WORKDIR /usr/src/app

RUN mkdir cert

# Copy only the manifest and lockfile so dependency install is cached
# and is not invalidated by unrelated source changes.
COPY package.json package-lock.json ./

# Install app dependencies straight from the lockfile.
# npm ci is reproducible and skips resolution entirely; it fails loudly if
# package.json and package-lock.json ever drift apart.
RUN npm ci

# Copy the rest of the application code to the working directory
COPY . .

# Build the app
RUN npm run build

# Expose the port that the app will run on
EXPOSE 3000

RUN chmod +x ./run.sh

# Define the command to run the app
CMD ["./run.sh"]