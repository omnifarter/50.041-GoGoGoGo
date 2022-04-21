
function NodeDisplay(props) {
    return (
        <table>
            <tbody>
                {props.keyStructure && Object.keys(props.keyStructure).map((key)=>(
                    <tr key={key}>
                        <td style={{width:'64px',height:'32px', fontWeight:'bold'}}>
                            Node {key}
                        </td>
                        {props.keyStructure[key].map((item)=>(
                        <td style={{width:"32px",height:'32px',textAlign:"center"}} key={item}>
                            {item}
                        </td>
                    ))}
                </tr>
                ))}
            </tbody>
    </table>
);
}

export default NodeDisplay;