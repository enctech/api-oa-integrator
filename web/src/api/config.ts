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
  if (id === "new") return {} as typeof sampleDetails;
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

export const createOAConfig = async (req: UpdateOAConfigRequest) => {
  return axios
    .post(`/config/snb-config`, {
      devices: req?.devices,
      facilities: req?.facilities,
      endpoint: req.endpoint,
      name: req.name,
      username: req.username,
      password: req.password,
    })
    .then((response) => response.data as typeof sampleDetails);
};

export const deleteOAConfig = async (id: string) => {
  return axios.delete(`/config/snb-config/${id}`);
};

export type SurchargeType = "percentage" | "exact";

// One config covers several sites: `name` is the config identity and is
// uniquely indexed, so four sites all named TNG cannot be four rows. Each
// site is a group owning the providerId S&B expects and the vendor's clientId.
export interface PlazaGroup {
  providerId: number;
  clientId: string;
  // Per site, not per facility: every facility on a site reports to the same
  // vendor location.
  vendorLocationId: string;
  // OA facility IDs belonging to this site.
  facilities: string[];
}

export interface IntegratorConfigs {
  id?: string;
  serviceProviderId: string;
  name: string;
  displayName: string;
  integratorName?: string;
  url: string;
  insecureSkipVerify: boolean;
  groups: PlazaGroup[];
  extra: Map<string, string>;
  taxRate: number;
  surcharge: number;
  surchargeType: SurchargeType;
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
  if (id === "new") return {} as IntegratorConfigs;
  return axios
    .get(`/config/integrator-config/${id}`)
    .then((response) => response.data as IntegratorConfigs);
};

export const updateIntegratorConfig = async (arg: IntegratorConfigs) => {
  console.log("updateIntegratorConfig", arg);
  const data = { ...arg };
  delete data["id"];
  return axios
    .put(`/config/integrator-config/${arg.id}`, {
      serviceProviderId: data.serviceProviderId,
      name: data.name,
      displayName: data.displayName,
      integratorName: data.integratorName,
      url: data.url,
      insecureSkipVerify: data.insecureSkipVerify,
      groups: data.groups,
      extra: JSON.parse(JSON.stringify(Object.fromEntries(data.extra))),
      taxRate: +arg.taxRate,
      surcharge: +arg.surcharge,
      surchargeType: arg.surchargeType,
    })
    .then((response) => response.data as IntegratorConfigs);
};

export const createIntegratorConfig = async (arg: IntegratorConfigs) => {
  const data = { ...arg };
  delete data["id"];
  return axios
    .post(`/config/integrator-config`, {
      serviceProviderId: data.serviceProviderId,
      name: data.name,
      displayName: data.displayName,
      integratorName: data.integratorName,
      url: data.url,
      insecureSkipVerify: data.insecureSkipVerify,
      groups: data.groups,
      extra: JSON.parse(JSON.stringify(Object.fromEntries(data.extra))),
      taxRate: +arg.taxRate,
      surcharge: +arg.surcharge,
      surchargeType: arg.surchargeType,
    })
    .then((response) => response.data as IntegratorConfigs);
};

export const deleteIntegratorConfig = async (id: string) => {
  return axios.delete(`/config/integrator-config/${id}`);
};
