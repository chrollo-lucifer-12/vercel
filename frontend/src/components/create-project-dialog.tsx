"use client";

import { ReactNode } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "./ui/dialog";
import AuthForm from "./custom/auth-form";
import { IdentificationBadgeIcon, LinkIcon } from "@phosphor-icons/react";
import { useCreateProjectMutation } from "@/hooks/use-project";

const CreateProject = ({ children }: { children: ReactNode }) => {
  const createProjectMutation = useCreateProjectMutation();

  return (
    <Dialog defaultOpen={false}>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create Project</DialogTitle>
          <DialogDescription>
            Create a new project with an existing git repository.
          </DialogDescription>
        </DialogHeader>
        <AuthForm
          hideBorder={true}
          loadingText="Creating Project..."
          onSubmit={(formData) => createProjectMutation.mutate({ formData })}
          submitText="Create Project"
          errors={createProjectMutation.data?.error}
          isPending={createProjectMutation.isPending}
          fields={[
            {
              id: "name",
              placeholder: "New Project",
              text: "Project Name",
              type: "text",
              icon: IdentificationBadgeIcon,
              required: true,
            },
            {
              id: "git_url",
              placeholder: "",
              text: "Repository url",
              type: "url",
              icon: LinkIcon,
            },
          ]}
        />
      </DialogContent>
    </Dialog>
  );
};

export default CreateProject;
