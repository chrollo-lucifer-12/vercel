"use client";

import { Deployment } from "@/lib/types";
import { ColumnDef } from "@tanstack/react-table";

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
];
