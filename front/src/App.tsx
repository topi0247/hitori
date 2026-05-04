import { useThemes } from "@/hooks/useThemes";

const App = () => {
  const { data, isLoading, isError } = useThemes();

  if (isLoading) return <p>loading...</p>;
  if (isError) return <p>error</p>;

  return (
    <ul>
      {data?.themes.map((theme) => (
        <li key={theme.id}>{theme.title}</li>
      ))}
    </ul>
  );
};

export default App;
