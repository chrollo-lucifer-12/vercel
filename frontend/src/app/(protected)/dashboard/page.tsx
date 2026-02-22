"use client";

import Navbar from "@/components/dashboard-page/navbar";
import SearchBar from "@/components/dashboard-page/search-bar";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "@/components/ui/input-group";
import { useSession } from "@/hooks/use-auth";
import { FileSearchIcon } from "@phosphor-icons/react";
import { useRouter } from "next/navigation";

const Dashboard = () => {
  const { data, isLoading } = useSession();
  console.log(data);
  const router = useRouter();
  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (!data) {
    router.replace("/signin");
    return null;
  }
  return (
    <div className="w-full p-4">
      <Navbar />
      <SearchBar />
    </div>
  );
};
export default Dashboard;
