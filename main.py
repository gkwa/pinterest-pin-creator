import csv
import re
import tempfile
import os
import argparse


def parse_arguments():
    parser = argparse.ArgumentParser(description="Process and modify CSV files.")
    parser.add_argument("input_file", help="Path to the input CSV file")
    parser.add_argument(
        "--increment-board", action="store_true", help="Increment the board number"
    )
    return parser.parse_args()


def clean_csv(file_path):
    temp_file = tempfile.NamedTemporaryFile(mode="w+", delete=False, newline="")

    with open(file_path, "r") as infile, temp_file:
        # Read the header
        header = infile.readline().strip()
        fieldnames = header.split(";")
        expected_field_count = len(fieldnames)

        # Write the header to the temp file
        temp_file.write(header + "\n")

        # Process each line
        for line in infile:
            if len(line.strip().split(";")) == expected_field_count:
                temp_file.write(line)

    # Replace the original file with the cleaned file
    temp_file.close()
    os.replace(temp_file.name, file_path)


def modify_csv(file_path, increment_board=False):
    temp_file = tempfile.NamedTemporaryFile(mode="w+", delete=False, newline="")

    with open(file_path, "r") as infile, temp_file:
        # Read the first line to get the fieldnames
        first_line = infile.readline().strip()
        fieldnames = first_line.split(";")

        reader = csv.DictReader(infile, fieldnames=fieldnames, delimiter=";")
        writer = csv.DictWriter(temp_file, fieldnames=fieldnames, delimiter=";")

        # Write the header
        writer.writeheader()

        # Process and write each row
        for row in reader:
            # Change 'created' to 'false'
            row["created"] = "false"

            # Increment the number in 'board' if the flag is set
            if increment_board:
                board = row["board"]
                match = re.search(r"(\D+)(\d+)", board)
                if match:
                    prefix, number = match.groups()
                    new_number = int(number) + 1
                    row["board"] = f"{prefix}{new_number}"

            writer.writerow(row)

    # Replace the original file with the modified file
    temp_file.close()
    os.replace(temp_file.name, file_path)


def main():
    args = parse_arguments()
    clean_csv(args.input_file)  # First, clean the CSV
    modify_csv(args.input_file, args.increment_board)  # Then, modify it


if __name__ == "__main__":
    main()
