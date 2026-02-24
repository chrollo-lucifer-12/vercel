"use client";

import { useProfile, useSession, useSignout } from "@/hooks/use-auth";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";
import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar";
import { SignOutIcon } from "@phosphor-icons/react";
import { Spinner } from "./ui/spinner";
import { Button } from "./ui/button";

const UserAvatar = () => {
  const { data, isLoading } = useProfile();
  const { mutate, isPending } = useSignout();

  const { data: session } = useSession();
  if (!session) return null;
  return (
    <DropdownMenu>
      <DropdownMenuTrigger>
        <Avatar className="h-8 w-8 rounded-lg">
          <AvatarImage src={"https://github.com/shadcn.png"} alt={data?.name} />
          <AvatarFallback className="rounded-lg">CN</AvatarFallback>
        </Avatar>
      </DropdownMenuTrigger>
      <DropdownMenuContent
        className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
        side={"right"}
        align="end"
        sideOffset={4}
      >
        <DropdownMenuLabel className="p-0 font-normal">
          <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
            <Avatar className="h-8 w-8 rounded-lg">
              <AvatarImage
                src={"https://github.com/shadcn.png"}
                alt={data?.name}
              />
              <AvatarFallback className="rounded-lg">CN</AvatarFallback>
            </Avatar>
            <div className="grid flex-1 text-left text-sm leading-tight">
              <span className="truncate font-medium">
                {isLoading ? "user" : data?.name}
              </span>
              <span className="truncate text-xs">
                {isLoading ? "user@gmail.com" : data?.email}
              </span>
            </div>
          </div>
        </DropdownMenuLabel>

        <DropdownMenuSeparator />
        <DropdownMenuItem asChild>
          <Button
            className="w-full"
            onClick={(e) => {
              e.preventDefault();
              mutate();
            }}
            disabled={isLoading || isPending}
          >
            {isPending ? (
              <>
                <Spinner className="animate-spin" />
                Logging out...
              </>
            ) : (
              <>
                <SignOutIcon />
                Log out
              </>
            )}
          </Button>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default UserAvatar;
