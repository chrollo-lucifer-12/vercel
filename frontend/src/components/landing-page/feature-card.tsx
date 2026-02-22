import { Icon } from "@phosphor-icons/react";
import { ReactNode } from "react";

interface FeatureCardProps {
  icon: Icon;
  title: string;
  content: ReactNode;
}

const FeatureCard = ({ content, icon: Icon, title }: FeatureCardProps) => {
  return (
    <div className="flex flex-col gap-y-1 max-w-[250px] ml-2 mb-2">
      <Icon className="text-purple-500" />
      <div>
        <p className="text-xs text-gray-700">
          <strong className="text-black">{title}</strong> {content}
        </p>
      </div>
    </div>
  );
};

export default FeatureCard;
