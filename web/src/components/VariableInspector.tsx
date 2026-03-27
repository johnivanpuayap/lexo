export interface VariableInfo {
  name: string;
  type: string;
  value: string;
}

interface VariableInspectorProps {
  variables: VariableInfo[];
}

export function VariableInspector({ variables }: VariableInspectorProps) {
  return (
    <div className="variable-inspector">
      <div className="panel-header">Variables</div>
      <div className="variable-list">
        {variables.length === 0 ? (
          <div className="empty-state">
            Run your program to see variables here.
          </div>
        ) : (
          <table className="variable-table">
            <thead>
              <tr>
                <th>Name</th>
                <th>Type</th>
                <th>Value</th>
              </tr>
            </thead>
            <tbody>
              {variables.map((v, i) => (
                <tr key={i}>
                  <td className="var-name">{v.name}</td>
                  <td>
                    <span className="type-badge">{v.type}</span>
                  </td>
                  <td className="var-value">{v.value}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
