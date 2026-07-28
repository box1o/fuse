import { Browser } from "./components/browser/browser";

const Main = () => {
    return (
        <div className="h-full w-full overflow-y-auto">
            <div className="grid min-h-full grid-cols-1 gap-6 p-6 lg:grid-cols-[minmax(0,1fr)_360px]">
                <main className="min-w-0">
                    <Browser />


                    <div className="mt-4 min-h-[1000px] rounded-xl bg-neutral-900" />
                </main>

                <aside className="space-y-4">
                    {Array.from({ length: 8 }).map((_, index) => (
                        <div
                            key={index}
                            className="grid grid-cols-[160px_1fr] gap-3"
                        >
                            <div className="aspect-video rounded-lg bg-neutral-800" />

                            <div className="space-y-2">
                                <div className="h-4 rounded bg-neutral-800" />
                                <div className="h-4 w-3/4 rounded bg-neutral-800" />
                                <div className="h-3 w-1/2 rounded bg-neutral-800" />
                            </div>
                        </div>
                    ))}
                </aside>
            </div>
        </div>
    );
};

export default Main;
