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
  console.log(response.data);
  sessionStorage.setItem("userData", JSON.stringify(response.data));
  return response.data;
};

export const logout = () => {
  sessionStorage.removeItem("userData");
};
