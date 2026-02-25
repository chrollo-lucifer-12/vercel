import Navbar from "@/components/dashboard-page/navbar";
import { ReactNode } from "react";

const Layout = ({ children }: { children: ReactNode }) => {
  return (
    <main className="w-full p-4">
      <Navbar />
      {children}
    </main>
  );
};

export default Layout;
