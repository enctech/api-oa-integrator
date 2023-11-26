import React, { useEffect, useState } from "react";
import {
  Autocomplete,
  Button,
  Chip,
  Container,
  TextField,
  Typography,
} from "@mui/material";
import { useMutation, useQuery } from "react-query";
import { getOAConfig, updateOAConfig } from "../../api/config";
import { useParams } from "react-router-dom";
import { getOAHealth } from "../../api/health";

const OaConfigsDetailsPage = () => {
  let { id } = useParams();
  const { data } = useQuery(["getOAConfig"], () => getOAConfig(id || ""));
  const { mutate, data: newData } = useMutation(
    "updateOAConfig",
    updateOAConfig,
  );

  const [isEditing, setIsEditing] = useState(false);
  const [name, setName] = useState(data?.name);
  const [facilities, setFacilities] = useState(data?.facilities);
  const [devices, setDevices] = useState(data?.devices);
  const [endpoint, setEndpoint] = useState(data?.endpoint);
  const [username, setUsername] = useState(data?.username);
  const [password, setPassword] = useState(data?.password);

  const { data: oaHealth } = useQuery("getOAHealth", () =>
    devices && devices.length > 0 && facilities && facilities.length > 0
      ? getOAHealth({
          device: devices[0],
          facility: facilities[0],
        })
      : Promise.resolve(null),
  );

  const handleSubmit = () => {
    mutate({
      id: id || "",
      name,
      facilities: facilities,
      devices: devices,
      endpoint,
      username,
      password,
    });
  };

  const reset = () => {
    if (!data) return;
    setName((newData || data).name);
    setFacilities((newData || data).facilities);
    setDevices((newData || data).devices);
    setEndpoint((newData || data).endpoint);
    setUsername((newData || data).username);
    setPassword((newData || data).password);
  };

  useEffect(reset, [data]);

  return (
    <Container className="p-16">
      <div className="flex content-between items-center justify-center mb-6">
        {oaHealth?.oa === "up" ? (
          <div
            className="w-5 h-5 mr-6 rounded-full
                inline-flex items-center justify-center
                bg-green-500"
          ></div>
        ) : (
          <div
            className="w-5 h-5 mr-6 rounded-full
                inline-flex items-center justify-center
                bg-red-500"
          ></div>
        )}
        <Typography variant="h5" component="h2">
          SnB Config Details
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
      <form>
        <div className="mb-8">
          <div>Name</div>
          <TextField
            fullWidth={true}
            variant="outlined"
            value={name}
            disabled={!isEditing}
            sx={{
              "& .MuiInputBase-input.Mui-disabled": {
                WebkitTextFillColor: "#000000",
              },
            }}
            onChange={(e) => setName(e.target.value)}
          />
        </div>
        {facilities && facilities.length > 0 && (
          <div className="mb-8">
            <div>Facilities</div>
            <Autocomplete
              clearIcon={false}
              options={[]}
              freeSolo
              multiple
              value={facilities}
              disabled={!isEditing}
              sx={{
                "& .Mui-disabled": {
                  opacity: 1,
                },
                "& .MuiButtonBase-root-MuiChip-root.Mui-disabled": {
                  opacity: 1,
                },
              }}
              renderTags={(value, props) =>
                value.map((option, index) => (
                  <Chip label={option} {...props({ index })} />
                ))
              }
              onChange={(e, value) => {
                setFacilities(value);
              }}
              renderInput={(params) => (
                <TextField
                  variant="outlined"
                  sx={{
                    "& .MuiInputBase-input.Mui-disabled": {
                      WebkitTextFillColor: "#000000",
                    },
                  }}
                  {...params}
                />
              )}
            />
          </div>
        )}
        {devices && devices.length > 0 && (
          <div className="mb-8">
            <div>Devices</div>
            <Autocomplete
              clearIcon={false}
              options={[]}
              freeSolo
              multiple
              value={devices}
              disabled={!isEditing}
              sx={{
                "& .Mui-disabled": {
                  opacity: 1,
                },
                "& .MuiButtonBase-root-MuiChip-root.Mui-disabled": {
                  opacity: 1,
                },
              }}
              renderTags={(value, props) =>
                value.map((option, index) => (
                  <Chip label={option} {...props({ index })} />
                ))
              }
              onChange={(e, value) => {
                setDevices(value);
              }}
              renderInput={(params) => (
                <TextField
                  variant="outlined"
                  sx={{
                    "& .MuiInputBase-input.Mui-disabled": {
                      WebkitTextFillColor: "#000000",
                    },
                  }}
                  {...params}
                />
              )}
            />
          </div>
        )}
        <div className="mb-8">
          <div>Endpoint</div>
          <TextField
            fullWidth={true}
            variant="outlined"
            value={endpoint}
            disabled={!isEditing}
            sx={{
              "& .MuiInputBase-input.Mui-disabled": {
                WebkitTextFillColor: "#000000",
              },
            }}
            onChange={(e) => setEndpoint(e.target.value)}
          />
        </div>
        <div className="mb-8">
          <div>Username</div>
          <TextField
            fullWidth={true}
            variant="outlined"
            value={username}
            disabled={!isEditing}
            sx={{
              "& .MuiInputBase-input.Mui-disabled": {
                WebkitTextFillColor: "#000000",
              },
            }}
            onChange={(e) => setUsername(e.target.value)}
          />
        </div>
        <div className="mb-8">
          <div>Password</div>
          <TextField
            fullWidth={true}
            variant="outlined"
            type="password"
            value={password}
            disabled={!isEditing}
            sx={{
              "& .MuiInputBase-input.Mui-disabled": {
                WebkitTextFillColor: "#000000",
              },
            }}
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>
        {isEditing && (
          <Button
            variant="contained"
            color="primary"
            fullWidth
            onClick={handleSubmit}
          >
            Save
          </Button>
        )}
      </form>
    </Container>
  );
};

export default OaConfigsDetailsPage;
