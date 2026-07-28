import { Browser } from "./components/browser/browser";
import { MachineListSkeleton } from "./machines";

const Main = () => {
    return (
        <div className="h-full w-full overflow-y-auto">
            <div className="grid min-h-full grid-cols-1 gap-6 p-6 lg:grid-cols-[minmax(0,1fr)_360px]">
                <main className="min-w-0">
                    <Browser />


                    <div className="mt-4 min-h-[1000px] rounded-xl bg-neutral-900" />
                </main>

                <MachineListSkeleton/>
            </div>
        </div>
    );
};

export default Main;
