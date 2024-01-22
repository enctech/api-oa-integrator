# Use an official Node runtime as a parent image
FROM node:20

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
CMD ["./run.sh"]