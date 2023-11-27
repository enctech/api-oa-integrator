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
  integratorName?: string;
  url: string;
  insecureSkipVerify: boolean;
  plazaIdMap: Map<string, string>;
}

export const getIntegratorConfigs = async () => {
  return axios
    .get(`/config/integrator-config`)
    .then((response) => response.data as IntegratorConfigs[]);
};

export const getIntegrators = async () => {
  return axios
    .get(`/config/integrators`)
    .then((response) => response.data as string[]);
};

export const getIntegratorConfig = async (id: string) => {
  if (id == "new") return {} as IntegratorConfigs;
  return axios
    .get(`/config/integrator-config/${id}`)
    .then((response) => response.data as IntegratorConfigs);
};

export const updateIntegratorConfig = async (arg: IntegratorConfigs) => {
  const data = { ...arg };
  delete data["id"];
  console.log(data);
  return axios
    .put(`/config/integrator-config/${arg.id}`, {
      clientId: data.clientId,
      providerId: data.providerId,
      serviceProviderId: data.serviceProviderId,
      name: data.name,
      integratorName: data.integratorName,
      url: data.url,
      insecureSkipVerify: data.insecureSkipVerify,
      plazaIdMap: JSON.parse(
        JSON.stringify(Object.fromEntries(arg.plazaIdMap)),
      ),
    })
    .then((response) => response.data as IntegratorConfigs);
};

export const createIntegratorConfig = async (arg: IntegratorConfigs) => {
  const data = { ...arg };
  delete data["id"];
  console.log(data);
  return axios
    .post(`/config/integrator-config`, {
      clientId: data.clientId,
      providerId: data.providerId,
      serviceProviderId: data.serviceProviderId,
      name: data.name,
      integratorName: data.integratorName,
      url: data.url,
      insecureSkipVerify: data.insecureSkipVerify,
      plazaIdMap: JSON.parse(
        JSON.stringify(Object.fromEntries(arg.plazaIdMap)),
      ),
    })
    .then((response) => response.data as IntegratorConfigs);
};
