# Use an official Node runtime as a parent image
FROM node:20

# Links the published package to this repository, which is what lets
# GITHUB_TOKEN publish it under the repo owner's namespace.
LABEL org.opencontainers.image.source=https://github.com/enctech/api-oa-integrator

# Set the working directory in the container
WORKDIR /usr/src/app

RUN mkdir cert

# Copy package.json and package-lock.json to the working directory
COPY package*.json ./

RUN npm install -g pnpm

# Install app dependencies
RUN pnpm install -d

# Copy the rest of the application code to the working directory
COPY . .

# Build the app
RUN pnpm run build

# Expose the port that the app will run on
EXPOSE 3000

RUN chmod +x ./run.sh

# Define the command to run the app
CMD ["npx", "serve", "build"]