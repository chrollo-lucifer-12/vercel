import { Project } from "@/lib/types";
import { Card, CardDescription, CardHeader, CardTitle } from "./ui/card";
import { Button } from "./ui/button";
import { GithubLogoIcon, TrashIcon } from "@phosphor-icons/react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "./ui/alert-dialog";
import { useDeleteProjectMutation } from "@/hooks/use-project";
import { Spinner } from "./ui/spinner";

const ProjectCard = ({ project, name }: { project: Project; name: string }) => {
  const router = useRouter();
  const { mutate, isPending } = useDeleteProjectMutation(name);
  return (
    <Card
      onClick={() => {
        router.push(`/project/${project.sub_domain}`);
      }}
      key={project.id}
      className={`cursor-pointer ${isPending ? "border-orange-500" : ""} hover:bg-[#E6F7F5] hover:text-[#26b3a6] border-2 border-transparent hover:border-[#1A8A81] transition duration-150`}
    >
      <CardHeader>
        <div className="flex flex-row items-center gap-2 justify-between">
          <Button variant={"outline"} className="rounded-full">
            <GithubLogoIcon />
          </Button>
          <div>
            <CardTitle>{project.name}</CardTitle>
            <CardDescription className="flex flex-col">
              <p> {project.sub_domain} </p>
              <Link
                className="hover:underline hover:text-[#26b3a6]"
                href={project.git_url}
              >
                {project.git_url}
              </Link>
            </CardDescription>
          </div>
          <AlertDialog>
            <AlertDialogTrigger asChild>
              <Button variant={"destructive"} disabled={isPending}>
                {!isPending ? <TrashIcon size={10} /> : <Spinner />}
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
                <AlertDialogDescription>
                  This action cannot be undone. This will permanently delete
                  your project from our servers.
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel>Cancel</AlertDialogCancel>
                <AlertDialogAction
                  onClick={() => {
                    mutate({ projectId: project.id });
                  }}
                  variant={"destructive"}
                >
                  Continue
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </div>
      </CardHeader>
    </Card>
  );
};

export default ProjectCard;
