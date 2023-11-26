import axios from "./axios";

export interface OAConfigResponse {
  id: string;
  name: string;
  endpoint: string;
  facilities: string[];
  devices: string[];
}

export const sampleDetails = {
  name: "KLCC Test",
  username: "6",
  password: "4711",
  endpoint: "https://helpdeskenctech.thruhere.net:8443",
  facilities: ["1230"],
  devices: ["101", "201"],
};

export const getOAConfigs = async () => {
  return axios
    .get(`/config/snb-config`)
    .then((response) => response.data as OAConfigResponse[]);
};

export const getOAConfig = async (id: string) => {
  return axios
    .get(`/config/snb-config/${id}`)
    .then((response) => response.data as typeof sampleDetails);
};

interface UpdateOAConfigRequest {
  id?: string;
  name?: string;
  username?: string;
  password?: string;
  devices?: string[];
  facilities?: string[];
  endpoint?: string;
}

export const updateOAConfig = async (req: UpdateOAConfigRequest) => {
  return axios
    .put(`/config/snb-config/${req.id}`, {
      devices: req?.devices,
      facilities: req?.facilities,
      endpoint: req.endpoint,
      name: req.name,
      username: req.username,
      password: req.password,
    })
    .then((response) => response.data as typeof sampleDetails);
};

export interface IntegratorConfigs {
  id?: string;
  clientId: string;
  providerId: number;
  serviceProviderId: string;
  name: string;
  url: string;
  insecureSkipVerify: boolean;
  plazaIdMap: {
    [key: string]: string;
  };
}

export const getIntegratorConfigs = async () => {
  return axios
    .get(`/config/integrator-config`)
    .then((response) => response.data as IntegratorConfigs[]);
};

export const getIntegratorConfig = async (id: string) => {
  return axios
    .get(`/config/integrator-config/${id}`)
    .then((response) => response.data as IntegratorConfigs);
};

export const updateIntegratorConfig = async (arg: IntegratorConfigs) => {
  const data = { ...arg };
  delete data["id"];
  return axios
    .put(`/config/integrator-config/${arg.id}`, data)
    .then((response) => response.data as IntegratorConfigs);
};
