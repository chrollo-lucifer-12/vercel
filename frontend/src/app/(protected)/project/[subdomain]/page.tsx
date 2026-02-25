import ProjectPage from "@/components/project-page";

const Page = async ({ params }: { params: Promise<{ subdomain: string }> }) => {
  const { subdomain } = await params;
  return <ProjectPage subdomain={subdomain} />;
};

export default Page;
