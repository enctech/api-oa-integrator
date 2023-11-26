import React, { useEffect, useState } from "react";
import { SubmitHandler, useFieldArray, useForm } from "react-hook-form";
import {
  Button,
  Container,
  TextField,
  Tooltip,
  tooltipClasses,
  TooltipProps,
} from "@mui/material";
import { useParams } from "react-router-dom";
import { useMutation, useQuery } from "react-query";
import {
  getIntegratorConfig,
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
  serviceProviderId: string;
  providerId: number;
  isInsecure: boolean;
  plazaIdMappers: { field1: string; field2: string }[];
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
    },
  });
  const { fields, append, remove } = useFieldArray({
    control,
    name: "plazaIdMappers",
  });

  let { id } = useParams();
  const { data } = useQuery(["getIntegratorConfig"], () =>
    getIntegratorConfig(id || ""),
  );
  const { mutate, data: newData } = useMutation(
    "updateIntegratorConfig",
    updateIntegratorConfig,
  );

  const [isEditing, setIsEditing] = useState(false);
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
    const map = new Map(Object.entries(data.plazaIdMap));
    map.forEach((value, key) => {
      if (fields.map((item) => item.field1).includes(key)) return;
      append({
        field1: key,
        field2: value as string,
      });
    });
  };

  useEffect(reset, [data]);

  const onSubmit: SubmitHandler<FormData> = (data) => {
    console.log(data);
    const plazaIdMap: { [key: string]: string } = {};
    data.plazaIdMappers.forEach((item) => {
      plazaIdMap[item.field1] = item.field2;
    });

    mutate({
      id: id || "",
      url: data.url,
      name: data.name,
      clientId: data.clientId,
      serviceProviderId: data.serviceProviderId,
      providerId: data.providerId,
      insecureSkipVerify: data.isInsecure,
      plazaIdMap: plazaIdMap,
    } satisfies IntegratorConfigs);
  };

  return (
    <Container className="p-16">
      <div className="flex content-between items-center justify-center mb-6">
        <Typography variant="h5" component="h2">
          {data?.name} Config Details
        </Typography>
        <div className="flex-1" />
        <Button
          variant="contained"
          color="primary"
          onClick={() => {
            setIsEditing(!isEditing);
            reset();
          }}
        >
          Edit
        </Button>
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
          <div className="mb-8">
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

        <div>
          <h2>Plaza ID Mapper</h2>
          {fields.map((field, index) => (
            <div key={field.id} className="flex">
              <div>
                <div>OA Facility ID</div>
                <TextField
                  {...register(`plazaIdMappers.${index}.field1` as const)}
                />
              </div>
              <div className="w-8" />
              <div>
                <div>Vendor Location ID</div>
                <TextField
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

const NoMaxWidthTooltip = styled(({ className, ...props }: TooltipProps) => (
  <Tooltip {...props} classes={{ popper: className }} />
))({
  [`& .${tooltipClasses.tooltip}`]: {
    maxWidth: "none",
  },
});

export default IntegratorConfigsDetails;
