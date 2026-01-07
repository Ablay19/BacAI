import { Solution } from "../types";
interface ProblemSolverProps {
    onSolve: (exercise: string, subject: string, level: string, language: string, mode: string, culturalContext: boolean) => Promise<Solution>;
    isLoading: boolean;
}
export default function ProblemSolver({ onSolve, isLoading, }: ProblemSolverProps): import("react/jsx-runtime").JSX.Element;
export {};
//# sourceMappingURL=ProblemSolver.d.ts.map