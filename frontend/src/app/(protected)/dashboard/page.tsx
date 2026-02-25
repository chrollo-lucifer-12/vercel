import AllProjectsClient from "@/components/dashboard-page/all-projects-client";
import Navbar from "@/components/dashboard-page/navbar";
import SearchBar from "@/components/dashboard-page/search-bar";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";

const Dashboard = async (props: {
  searchParams?: Promise<{ name?: string }>;
}) => {
  const cookieStore = await cookies();
  const refreshToken = cookieStore.get("refresh_token")?.value;

  if (!refreshToken) {
    redirect("/signin");
  }

  const searchParams = await props.searchParams;
  const name = searchParams?.name || "";

  return (
    <>
      <SearchBar />
      <AllProjectsClient name={name} />
    </>
  );
};
export default Dashboard;
