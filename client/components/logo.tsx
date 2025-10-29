import { FaHouseChimney } from "react-icons/fa6";

const Logo = () => {
  return (
    <div className="flex items-center justify-center gap-2 text-slate-900">
      <FaHouseChimney size={30} className="fill-primary" />
      <h1 className="text-3xl font-bold">CozyStay</h1>
    </div>
  );
};

export default Logo;
