import Navbar from "@/components/dashboard-page/navbar";
import SearchBar from "@/components/dashboard-page/search-bar";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";

const Dashboard = async () => {
  const cookieStore = await cookies();
  const refreshToken = cookieStore.get("refresh_token")?.name;

  if (!refreshToken) {
    redirect("/signin");
  }

  return (
    <div className="w-full p-4">
      <Navbar />
      <SearchBar />
    </div>
  );
};
export default Dashboard;
