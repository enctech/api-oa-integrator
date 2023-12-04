import React from "react";
import {
  Button,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
} from "@mui/material";
import { useMutation, useQuery } from "react-query";
import { deleteUser, users } from "../api/auth";
import DeleteIcon from "@mui/icons-material/Delete";
import AlertDialog from "../components/dialog";

const UsersPage = () => {
  const { data, refetch } = useQuery("users", () => users());
  const { mutate: mutateDelete } = useMutation("deleteUser", deleteUser, {
    onSuccess: () => refetch(),
    onSettled: () => setShowDeleteUserDialog([false, ""]),
  });

  const [showDeleteUserDialog, setShowDeleteUserDialog] = React.useState<
    [boolean, string]
  >([false, ""]);

  const createNewUser = () => {};

  return (
    <>
      <div>
        <Button variant="contained" color="primary" onClick={createNewUser}>
          New
        </Button>
        <TableContainer component={Paper} className="mt-4">
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell>Permissions</TableCell>
                <TableCell />
              </TableRow>
            </TableHead>
            <TableBody>
              {data?.map((row) => (
                <TableRow className={"cursor-pointer"} key={row.userId}>
                  <TableCell>{row.name}</TableCell>
                  <TableCell>{row.permissions.join(",")}</TableCell>
                  <TableCell
                    onClick={() => setShowDeleteUserDialog([true, row.userId])}
                    style={{
                      width: "10px",
                      whiteSpace: "normal",
                      wordWrap: "break-word",
                    }}
                  >
                    <DeleteIcon />
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </div>
      <AlertDialog
        isOpen={showDeleteUserDialog[0]}
        handleClose={() => setShowDeleteUserDialog([false, ""])}
        title={"Delete user"}
        description={"Are you sure you want to delete this user?"}
        buttons={[
          <Button
            key="cancel"
            onClick={() => setShowDeleteUserDialog([false, ""])}
            color="primary"
          >
            Cancel
          </Button>,
          <Button
            key="logout"
            onClick={() => {
              mutateDelete(showDeleteUserDialog[1]);
            }}
            color="primary"
            autoFocus
          >
            Delete
          </Button>,
        ]}
      />
    </>
  );
};

export default UsersPage;
