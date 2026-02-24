import Image from "next/image";
import Link from "next/link";
import UserAvatar from "../user-avatar";
import { Suspense } from "react";

const Navbar = () => {
  return (
    <div className="w-full flex justify-between items-center ">
      <Link href="/" className="flex items-center">
        <Image src="/logo.png" alt="logo" width={100} height={50} />
      </Link>
      <UserAvatar />
    </div>
  );
};

export default Navbar;
