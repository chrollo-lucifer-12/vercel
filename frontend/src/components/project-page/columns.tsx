"use client";

import { Deployment } from "@/lib/types";
import { ColumnDef } from "@tanstack/react-table";
import BuildDialog from "./build-dialog";
import { Button } from "../ui/button";
import { DotsThreeCircleIcon } from "@phosphor-icons/react";

export const columns: ColumnDef<Deployment>[] = [
  {
    accessorKey: "sequence",
    header: "Sequence",
  },
  {
    accessorKey: "created_at",
    header: "Created At",
  },
  {
    accessorKey: "status",
    header: "Status",
  },
  {
    accessorKey: "action",
    header: "Logs",
    cell: ({ row }) => {
      const id = row.original.id;
      return (
        <BuildDialog deploymentId={id}>
          <Button>
            <DotsThreeCircleIcon />
          </Button>
        </BuildDialog>
      );
    },
  },
];
