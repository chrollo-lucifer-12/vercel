import ProjectPage from "@/components/project-page";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";

const Page = async ({ params }: { params: Promise<{ subdomain: string }> }) => {
  const cookieStore = await cookies();
  const refreshToken = cookieStore.get("refresh_token")?.value;

  if (!refreshToken) {
    redirect("/signin");
  }
  const { subdomain } = await params;
  return <ProjectPage subdomain={subdomain} />;
};

export default Page;
