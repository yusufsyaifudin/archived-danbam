#!/usr/bin/env python

import sys, os
import time

__location__ = os.path.realpath(os.path.join(os.getcwd(), os.path.dirname(__file__)))

def snake_to_camel(word):
    return ''.join(x.capitalize() or '_' for x in word.split('_'))

def parse_template(current_time, migration_name):
    file = open(os.path.join(__location__, 'migration_template'), 'r')
    template = ''.join(file.readlines())
    file.close()

    return template \
        .replace("{struct_name}", snake_to_camel(migration_name) + str(current_time)) \
        .replace("{migration_name}", migration_name) \
        .replace("{current_time}", str(current_time))

def write_migration(migration_name):
    current_time = int(time.time())
    path = os.path.join(__location__, '../../assets/migrate/'
                        + str(current_time) + '_'
                        + migration_name + '.go')
    f = open(path, "w+")
    f.write(parse_template(current_time, migration_name))
    f.close()
    return

functions = {
    'write_migration': write_migration,
}

if __name__ == '__main__':
    func = functions[sys.argv[1]]
    args = sys.argv[2:]

    sys.exit(func(*args))