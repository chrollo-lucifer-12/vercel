import { FilePlusIcon, PlusIcon } from "@phosphor-icons/react";
import {
  Empty,
  EmptyContent,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from "./ui/empty";
import { Button } from "./ui/button";
import CreateProject from "./create-project-dialog";

const ProjectEmpty = () => {
  return (
    <Empty className="bg-muted/30 h-full w-full ">
      <EmptyHeader>
        <EmptyMedia variant="icon">
          <FilePlusIcon />
        </EmptyMedia>
        <EmptyTitle>No Projects</EmptyTitle>
        <EmptyDescription className="max-w-xs text-pretty">
          You haven't created any projects yet.
        </EmptyDescription>
      </EmptyHeader>
      <EmptyContent>
        <CreateProject>
          <Button size={"icon"}>
            <PlusIcon />
          </Button>
        </CreateProject>
      </EmptyContent>
    </Empty>
  );
};

export default ProjectEmpty;
