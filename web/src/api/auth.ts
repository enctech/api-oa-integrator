import axios from "./axios";

interface LoginRequest {
  username: string;
  password: string;
}

export const login = async ({ username, password }: LoginRequest) => {
  const response = await axios.post("/auth/login", {
    username,
    password,
  });
  return response.data;
};
