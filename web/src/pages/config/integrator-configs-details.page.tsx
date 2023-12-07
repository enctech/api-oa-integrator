import React, { useEffect, useState } from "react";
import { SubmitHandler, useFieldArray, useForm } from "react-hook-form";
import {
  Button,
  Checkbox,
  Container,
  FormControlLabel,
  Radio,
  RadioGroup,
  TextField,
  Tooltip,
  tooltipClasses,
  TooltipProps,
} from "@mui/material";
import { useParams } from "react-router-dom";
import { useMutation, useQuery } from "react-query";
import {
  createIntegratorConfig,
  getIntegratorConfig,
  getIntegrators,
  IntegratorConfigs,
  updateIntegratorConfig,
} from "../../api/config";
import Typography from "@mui/material/Typography";
import InfoIcon from "@mui/icons-material/Info";
import { styled } from "@mui/material/styles";

interface FormData {
  url: string;
  name: string;
  clientId: string;
  integratorName?: string;
  serviceProviderId: string;
  providerId: number;
  isInsecure: boolean;
  plazaIdMappers: { field1: string; field2: string }[];
  extra: any[];
}

const IntegratorConfigsDetails = () => {
  const {
    control,
    handleSubmit,
    register,
    setValue,
    formState: { errors },
  } = useForm<FormData>({
    defaultValues: {
      url: "",
      name: "",
      clientId: "",
      serviceProviderId: "",
      providerId: 0,
      isInsecure: false,
      plazaIdMappers: [],
      extra: [],
    },
  });
  let { id } = useParams();
  const { fields, remove, append, update } = useFieldArray({
    control,
    name: "plazaIdMappers",
  });

  const { data, isLoading } = useQuery(
    ["getIntegratorConfig", id],
    () => getIntegratorConfig(id || ""),
    {
      onSuccess: (data) => {
        if (data.integratorName) {
          setIntegrator(data.integratorName);
        }
      },
    },
  );
  const [integrator, setIntegrator] = useState<string>(
    data?.integratorName || "",
  );
  const { data: integrators } = useQuery(["getIntegrators"], () =>
    getIntegrators(),
  );
  const { mutate, data: updatedData } = useMutation(
    "updateIntegratorConfig",
    updateIntegratorConfig,
    {
      onSettled: () => {
        setIsEditing(false);
      },
    },
  );

  const { mutate: create, data: newData } = useMutation(
    "createIntegratorConfig",
    createIntegratorConfig,
    {
      onSettled: () => {
        setIsEditing(false);
      },
    },
  );

  const [isEditing, setIsEditing] = useState(id === "new");
  const [name, setName] = useState("");

  const reset = () => {
    if (!data) return;
    setValue("url", data.url);
    setValue("name", data.name);
    setName(data.name);
    setValue("clientId", data.clientId);
    setValue("providerId", data.providerId);
    setValue("isInsecure", data.insecureSkipVerify);
    setValue("serviceProviderId", data.serviceProviderId);
    setValue("integratorName", data.integratorName);
    if (data.plazaIdMap) {
      const keys = Object.keys(data.plazaIdMap);
      keys.forEach((key, index) => {
        update(index, {
          field1: key,
          field2: new Map(Object.entries(data.plazaIdMap)).get(key) || "",
        });
      });
    } else {
      update(0, { field1: "", field2: "" });
    }

    if (data.extra) {
      const extraDataForm = buildExtraDataForForm(
        data.integratorName || "",
        data.extra,
      );

      extraDataForm.forEach((value, key) => {
        setValue(key as any, value);
      });
    }
  };

  useEffect(reset, [data]);

  const onSubmit: SubmitHandler<FormData> = (data) => {
    console.log(data);
    const plazaIdMap: Map<string, string> = new Map();
    data.plazaIdMappers.forEach((item) => {
      plazaIdMap.set(item.field1, item.field2);
    });

    if (id == "new") {
      create({
        id: id || "",
        url: data.url,
        name: data.name,
        clientId: data.clientId,
        serviceProviderId: data.serviceProviderId,
        providerId: data.providerId,
        insecureSkipVerify: data.isInsecure,
        plazaIdMap: plazaIdMap,
        integratorName: data.integratorName,
        extra: buildExtraDataForVendor(data.integratorName || "", data.extra),
      } satisfies IntegratorConfigs);
      return;
    }

    mutate({
      id: id || "",
      url: data.url,
      name: data.name,
      clientId: data.clientId,
      serviceProviderId: data.serviceProviderId,
      providerId: data.providerId,
      insecureSkipVerify: data.isInsecure,
      plazaIdMap: plazaIdMap,
      integratorName: data.integratorName,
      extra: buildExtraDataForVendor(data.integratorName || "", data.extra),
    } satisfies IntegratorConfigs);
  };

  return (
    <Container className="p-16">
      <div className="flex content-between items-center justify-center mb-6">
        <Typography variant="h5" component="h2">
          {data?.name} Config Details
        </Typography>
        <div className="flex-1" />
        {id !== "new" && (
          <Button
            variant="contained"
            color="primary"
            onClick={() => {
              setIsEditing(!isEditing);
              reset();
            }}
          >
            {isEditing && "Cancel"} Edit
          </Button>
        )}
      </div>
      <form onSubmit={handleSubmit(onSubmit)}>
        <div>
          <div className="mb-8">
            <div>
              Name
              <NoMaxWidthTooltip
                title={`This will be identifier for SnB to use to call OA system. Example, http://localhost:8080/oa/${
                  name || "{name}"
                }/AuthorizationService3rdParty`}
              >
                <InfoIcon />
              </NoMaxWidthTooltip>
            </div>
            <TextField
              fullWidth={true}
              variant="outlined"
              disabled={!isEditing}
              sx={{
                "& .MuiInputBase-input.Mui-disabled": {
                  WebkitTextFillColor: "#000000",
                },
              }}
              {...register("name", {
                onChange: (e) => setName(e.target.value),
              })}
            />
          </div>
        </div>
        <div>
          <div className="flex">
            <div className="mb-8 flex-1">
              <div>URL</div>
              <TextField
                fullWidth={true}
                variant="outlined"
                disabled={!isEditing}
                sx={{
                  "& .MuiInputBase-input.Mui-disabled": {
                    WebkitTextFillColor: "#000000",
                  },
                }}
                {...register("url")}
              />
            </div>
            <div className="w-8" />
            <FormControlLabel
              {...register("isInsecure")}
              control={<Checkbox />}
              disabled={!isEditing}
              sx={{
                "& .MuiTypography-root": {
                  WebkitTextFillColor: "#000000",
                },
              }}
              label="Insecure endpoint"
            />
          </div>
        </div>
        <div>
          <div className="mb-8">
            <div>Provider ID (For OA)</div>
            <TextField
              fullWidth={true}
              variant="outlined"
              disabled={!isEditing}
              sx={{
                "& .MuiInputBase-input.Mui-disabled": {
                  WebkitTextFillColor: "#000000",
                },
              }}
              {...register("providerId")}
            />
          </div>
        </div>
        <div>
          <div className="mb-8">
            <div>Client ID (Defined by integrator)</div>
            <TextField
              fullWidth={true}
              variant="outlined"
              disabled={!isEditing}
              sx={{
                "& .MuiInputBase-input.Mui-disabled": {
                  WebkitTextFillColor: "#000000",
                },
              }}
              {...register("clientId")}
            />
          </div>
        </div>
        <div>
          <div className="mb-8">
            <div>
              Service Provider ID (Defined by integrator)
              <Tooltip
                className="ml-2"
                title="For any identifier integrator used to define OA system"
              >
                <InfoIcon />
              </Tooltip>
            </div>
            <TextField
              fullWidth={true}
              variant="outlined"
              disabled={!isEditing}
              sx={{
                "& .MuiInputBase-input.Mui-disabled": {
                  WebkitTextFillColor: "#000000",
                },
              }}
              {...register("serviceProviderId")}
            />
          </div>
        </div>

        {data && data?.integratorName ? (
          <RadioGroup
            defaultValue={data?.integratorName}
            className="mb-8"
            {...register(`integratorName` as const, {
              onChange: (e) => {
                setIntegrator(e.target.value);
              },
            })}
          >
            <div>
              Integrator type
              <Tooltip
                className="ml-2"
                title="Which integrator is this configuration system is for?"
              >
                <InfoIcon />
              </Tooltip>
            </div>
            {integrators?.map((value) => (
              <FormControlLabel
                disabled={!isEditing}
                key={value}
                value={value}
                control={<Radio />}
                label={value.toUpperCase()}
                sx={{
                  "& .MuiTypography-root": {
                    WebkitTextFillColor: "#000000",
                  },
                }}
              />
            ))}
          </RadioGroup>
        ) : (
          <RadioGroup
            key={"new"}
            defaultValue={integrators && integrators[0]}
            className="mb-8"
            {...register(`integratorName` as const, {
              onChange: (e) => {
                setIntegrator(e.target.value);
              },
            })}
          >
            <div>
              Integrator type
              <Tooltip
                className="ml-2"
                title="Which integrator is this configuration system is for?"
              >
                <InfoIcon />
              </Tooltip>
            </div>
            {integrators?.map((value) => (
              <FormControlLabel
                disabled={!isEditing}
                key={value}
                value={value}
                control={<Radio />}
                label={value.toUpperCase()}
                sx={{
                  "& .MuiTypography-root": {
                    WebkitTextFillColor: "#000000",
                  },
                }}
              />
            ))}
          </RadioGroup>
        )}
        {integrator === "tng" && (
          <div className="mb-8">
            <div>
              SSH Key
              <Tooltip
                className="ml-2"
                title="TNG Use this for server legibility"
              >
                <InfoIcon />
              </Tooltip>
            </div>
            <TextField
              fullWidth={true}
              variant="outlined"
              multiline={true}
              minRows={6}
              disabled={!isEditing}
              sx={{
                "& .MuiInputBase-input.Mui-disabled": {
                  WebkitTextFillColor: "#000000",
                },
                "& .MuiOutlinedInput-root": {
                  backgroundColor: isEditing ? "white" : "#d2d5d8",
                },
              }}
              {...register("extra.0")}
            />
          </div>
        )}

        <div>
          <h2>Plaza ID Mapper</h2>
          {data &&
            fields.map((field, index) => (
              <div
                key={`${field.id}-${field.field1}-${field.field2}`}
                className="flex"
              >
                <div>
                  <div>OA Facility ID</div>
                  <TextField
                    disabled={!isEditing}
                    sx={{
                      "& .MuiInputBase-input.Mui-disabled": {
                        WebkitTextFillColor: "#000000",
                      },
                    }}
                    {...register(`plazaIdMappers.${index}.field1` as const)}
                  />
                </div>
                <div className="w-8" />
                <div>
                  <div>Vendor Location ID</div>
                  <TextField
                    disabled={!isEditing}
                    sx={{
                      "& .MuiInputBase-input.Mui-disabled": {
                        WebkitTextFillColor: "#000000",
                      },
                    }}
                    {...register(`plazaIdMappers.${index}.field2` as const)}
                  />
                </div>
                <div className="w-8" />
                {isEditing && (
                  <Button type="button" onClick={() => remove(index)}>
                    Remove
                  </Button>
                )}
              </div>
            ))}
          {isEditing && (
            <Button
              type="button"
              onClick={() => append({ field1: "", field2: "" })}
            >
              Add Field
            </Button>
          )}
        </div>
        <div className="h-8" />

        {isEditing && (
          <Button variant="contained" color="primary" fullWidth type="submit">
            Save
          </Button>
        )}
      </form>
    </Container>
  );
};

const buildExtraDataForVendor = (
  vendor: string,
  input: any[],
): Map<string, string> => {
  const out = new Map<string, any>();
  if (vendor == "tng") {
    out.set("sshKey", input[0]);
  }
  return out;
};

const buildExtraDataForForm = (
  vendor: string,
  data: { [key: string]: any },
): Map<string, string> => {
  const out = new Map<string, any>();
  if (vendor == "tng") {
    out.set("extra.0", data["sshKey"]);
  }
  return out;
};

const NoMaxWidthTooltip = styled(({ className, ...props }: TooltipProps) => (
  <Tooltip {...props} classes={{ popper: className }} />
))({
  [`& .${tooltipClasses.tooltip}`]: {
    maxWidth: "none",
  },
});

export default IntegratorConfigsDetails;
