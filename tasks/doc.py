'''A series of tasks to generate documentation'''

from pathlib import Path

from invoke import task, Collection
from jinja2 import Environment, PackageLoader

from snooze.utils.config import *

#from tasks.utils import get_versions, print_github_kv, get_paths

def compute_type(prop: dict) -> str:
    if prop['type'] == 'array' and prop.get('items', {}).get('type'):
        return f"array[{prop['items']['type']}]"
    else:
        return prop['type']

def get_ref(obj: dict) -> str:
    if '$ref' in obj:
        return obj['$ref']
    if 'allOf' in obj:
        return obj['allOf'][0]['$ref']
    return None

def append_dot(line: str) -> str:
    if line[-1] != '.':
        line += '.'
    return line

def prop_to_markdown(name: str, prop: dict, required: bool) -> str:
    prop_line = '- '
    prop_line += f"`{name}`"
    if 'type' in prop:
        prop_line += f" ({compute_type(prop)})"
    if required:
        prop_line += ' <required>'
    else:
        ref = get_ref(prop)
        if ref:
            ref_name = ref.split('/')[-1]
            ref_link = ref
            prop_line += f" ([{ref_name}]({ref_link}))"
    prop_line += ':'
    if 'title' in prop and (
        prop['title'] != name.capitalize().replace('_', ' ')
        or 'description' not in prop
    ):
        prop_line += f" {append_dot(prop['title'])}"
    if 'description' in prop:
        prop_line += f" {append_dot(prop['description'])}"
    if 'env' in prop:
        prop_line += f" Environment variable: `{prop['env']}`."
    if 'default' in prop:
        prop_line += f" Default: `{prop['default']}`."
    return prop_line

def definition_to_markdown(name: str, definition: dict) -> str:
    output = ''
    if 'title' in definition and definition['title'] != name:
        output += f"### {definition['title']} ({name})\n\n"
    else:
        output += f"### {name}\n\n"
    if 'description' in definition:
        output += f"{definition['description']}\n\n"
    if 'properties' in definition:
        output += "#### Properties\n\n"
        for name, prop in definition['properties'].items():
            required = (name in definition.get('required', []))
            output += f"{prop_to_markdown(name, prop, required)}\n"
    return output

def schema_to_markdown(schema: dict) -> str:
    '''Convert the schema to markdown'''
    output = ''
    if 'title' in schema:
        output += f"# {schema['title']}\n\n"
    if 'description' in schema:
        output += f"{schema['description']}\n\n"
    if 'properties' in schema:
        output += "## Properties\n\n"
        for name, prop in schema['properties'].items():
            required = (name in schema.get('required', []))
            output += f"{prop_to_markdown(name, prop, required)}\n"
        output += "\n"
    if 'definitions' in schema:
        output += "## Definitions\n\n"
        for name, definition in schema['definitions'].items():
            output += f"{definition_to_markdown(name, definition)}\n"
        output += "\n"
    return output

@task
def config(ctx):
    '''Generate documentation for configuration files'''
    doc_path = Path('doc/config.md')
    doc_str = ''
    doc_str += schema_to_markdown(CoreConfig.schema())
    doc_str += schema_to_markdown(GeneralConfig.schema())
    doc_str += schema_to_markdown(HousekeeperConfig.schema())
    doc_str += schema_to_markdown(NotificationConfig.schema())
    doc_str += schema_to_markdown(LdapConfig.schema())
    doc_path.write_text(doc_str, encoding='utf-8')
    print(f"Documentation generated in {doc_path}")

ns = Collection('doc')
ns.add_task(config)
