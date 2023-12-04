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
  sessionStorage.setItem("userData", JSON.stringify(response.data));
  return response.data;
};

export const logout = () => {
  sessionStorage.removeItem("userData");
};

const usersResponse = [
  {
    username: "string",
    name: "TEST",
    userId: "23224a79-7ff0-4432-bee4-f918259e764d",
    permissions: ["admin"],
  },
];

export const users = async () => {
  return axios
    .get("/auth/users")
    .then((res) => res.data as typeof usersResponse);
};

const createUserRequest = {
  name: "string",
  password: "string",
  permissions: ["string"],
  username: "string",
};

export const createUser = async (req: typeof createUserRequest) => {
  return axios
    .post("/auth/user", req)
    .then((res) => res.data as typeof usersResponse);
};

export const deleteUser = async (id: string) => {
  return axios.delete(`/auth/user/${id}`);
};
